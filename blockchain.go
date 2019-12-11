package main

import (
	"bytes"
//	"crypto/ecdsa"
//	"encoding/hex"
	"errors"
	"fmt"
	"log"
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

func CreateBlockchain(address string) *Blockchain {
	// TODO: sort and open correct bc file when there are multiple
	if dbExists() {
		fmt.Println("--> another blockchain already exists ")
		os.Exit(1)
	}
	var tip []byte
  db, err := bolt.Open(dbFile, 0600, nil)
  if err != nil {
		panic(err)
	}
	fmt.Printf("--> Creating genesis block\n")
  err = db.Update(func(tx *bolt.Tx) error {
		bucket1, err := tx.CreateBucket([]byte(blockchainBucket))
		if err != nil {
			panic(err)
		}
		genesistransaction := NewGenesisTransaction(address, "BLOCKCHAIN")
		genesisblock := NewGenesisBlock(genesistransaction)
		err = bucket1.Put(genesisblock.Hash, genesisblock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = bucket1.Put([]byte("l"), genesisblock.Hash)
		if err != nil {
			panic(err)
		}
		tip = genesisblock.Hash
		_, err = tx.CreateBucket([]byte(miningBucket))
		if err != nil {
			panic(err)
		}
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc
}

func LoadBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockchainBucket))
		tip = bucket.Get([]byte("l"))
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc
}

func (bc *Blockchain) AddToQueue(tr *Transaction) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(miningBucket))
		err := bucket.Put(tr.ID, tr.Serialize())
		return err
	})
	if(err != nil){
		panic(err)
	}
}

func (bc *Blockchain) MineBlock(size int) {
	//Get has from prev block
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockchainBucket))
		lastHash = bucket.Get([]byte("l"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	//Get the last transactions
	txs := make([]*Transaction, 0)
	err = bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(miningBucket))
		cur := bucket.Cursor()
		for k, v := cur.First(); k != nil && len(txs) <= size; k, v = cur.Next() {
			tx := DeserializeToTX(v)
			if tx.Validate() {
				fmt.Printf("--> Transaction verified: %x\n", tx.ID)
				txs = append(txs, tx)
			}
		}
		return nil
	})
	if len(txs) == 0 {
		panic(errors.New("No transactions to mine"))
	}
	if err != nil {
	}
  block := NewBlock(txs, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockchainBucket))
		err := bucket.Put(block.Hash, block.Serialize())
		if err != nil {
			panic(err)
		}
		err = bucket.Put([]byte("l"), block.Hash)
		if err != nil {
			panic(err)
		}
		bc.tip = block.Hash
		return nil
	})
	if err != nil {
		panic(err)
	}
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
	return Transaction{}, errors.New("Not found")
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
