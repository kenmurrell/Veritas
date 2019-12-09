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
const bucketName = "blocks"


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
		log.Panic(err)
	}
  err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			log.Panic(err)
		}
		genesistransaction := NewGenesisTransaction(address, "BLOCKCHAIN")
		genesisblock := NewGenesisBlock(genesistransaction)
		err = bucket.Put(genesisblock.Hash, genesisblock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte("l"), genesisblock.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesisblock.Hash
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
		bucket := tx.Bucket([]byte(bucketName))
		tip = bucket.Get([]byte("l"))
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		lastHash = bucket.Get([]byte("l"))
		return nil
	})
	if err != nil {
		panic(err)
	}
  block := NewBlock(transactions, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
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
		bucket := tx.Bucket([]byte(bucketName))
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
