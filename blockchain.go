package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blockchainBucket = "blocks"
const miningBucket = "queue"


type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateBlockchain(file, address string) (*Blockchain, error) {
	if dbExists() {
		return nil, errors.New("ERROR: Blockchain file already exists")
	}
	var tip []byte
  db, err := bolt.Open(file, 0600, nil)
  if err != nil {
		return nil, errors.New("ERROR: Cannot open blockchain file as database")
	}
  err = db.Update(func(tx *bolt.Tx) error {
		bucket1, err := tx.CreateBucket([]byte(blockchainBucket))
		if err != nil {
			return err
		}
		genesistransaction := NewGenesisTransaction(address, "BLOCKCHAIN")
		genesisblock := NewGenesisBlock(genesistransaction)
		err = bucket1.Put(genesisblock.Hash, genesisblock.Serialize())
		if err != nil {
			return err
		}
		err = bucket1.Put([]byte("l"), genesisblock.Hash)
		if err != nil {
			return err
		}
		tip = genesisblock.Hash
		_, err = tx.CreateBucket([]byte(miningBucket))
		if err != nil {
			return err
		}
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc, err
}

func LoadBlockchain(file string) (*Blockchain, error) {
	var tip []byte
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return nil, errors.New("ERROR: Cannot open blockchain file as database\n")
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockchainBucket))
		tip = bucket.Get([]byte("l"))
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc, err
}

func (bc *Blockchain) AddToQueue(tr *Transaction) error {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(miningBucket))
		err := bucket.Put(tr.ID, tr.Serialize())
		return err
	})
	return err
}

func (bc *Blockchain) MineBlock(size int) error {
	//Get hash from prev block
	var lastHash []byte
	_ = bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockchainBucket))
		lastHash = bucket.Get([]byte("l"))
		return nil
	})
	if lastHash == nil {
		return errors.New("ERROR: cannot link to previous hash\n")
	}
	//Get the last transactions
	txs := make([]*Transaction, 0)
	_ = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(miningBucket))
		cur := bucket.Cursor()
		for k, v := cur.First(); k != nil && len(txs) < size; k, v = cur.Next() {
			tx := DeserializeToTX(v)
			if tx.Validate() {
				fmt.Printf("--> Transaction verified: %x\n", tx.ID)
				txs = append(txs, tx)
			}
			_ = bucket.Delete(k)
		}
		return nil
	})
	if len(txs) == 0 {
		return errors.New("ERROR: No transactions to mine\n")
	}
  block := NewBlock(txs, lastHash)
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockchainBucket))
		err := bucket.Put(block.Hash, block.Serialize())
		if err != nil {
			return err
		}
		err = bucket.Put([]byte("l"), block.Hash)
		if err != nil {
			return err
		}
		bc.tip = block.Hash
		return nil
	})
	return err
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	_ = bci.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockchainBucket))
		encodedBlock := bucket.Get(bci.currentHash)
		block = Deserialize(encodedBlock)
		return nil
	})
	bci.currentHash = block.PrevHash
	return block
}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()
	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("ERROR: Transaction not found\n")
}

func (bc *Blockchain) CalculateBalance(address string) int {
	bci := bc.Iterator()
	balance := 0
	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if strings.Compare(address, string(tx.ToAddr)) == 0 {
				balance = balance + tx.Amount
			} else if strings.Compare(address, tx.FromAddr) == 0 {
				balance = balance - tx.Amount
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return balance
}

func (bc *Blockchain) Close() {
	bc.db.Close()
}
