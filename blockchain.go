package main

import (
	"bytes"
//	"crypto/ecdsa"
//	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const bucket = "blocks"


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
	var temp []byte
  db, err := bolt.Open(dbFile, 0600, nil)
  if err != nil {
		log.Panic(err)
	}
	defer db.Close()
  err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(bucket))
		if err != nil {
			log.Panic(err)
		}
		genesistransaction := NewGenesisTransaction(address, "XXXXX")
		genesisblock := NewGenesisBlock(genesistransaction)
		err = bucket.Put(genesisblock.Hash, genesisblock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte("l"), genesisblock.Hash)
		if err != nil {
			log.Panic(err)
		}
		temp = genesisblock.Hash
		return nil
	})
	bc := Blockchain{temp, db}
	return &bc
}

func LoadBlockchain(address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		tip = b.Get([]byte("l"))
		return nil
	})
	bc := Blockchain{tip, db}
	return &bc
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		panic(err)
	}
  block := NewBlock(transactions, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(block.Hash, block.Serialize())
		if err != nil {
			panic(err)
		}
		err = b.Put([]byte("l"), block.Hash)
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

func (i *BlockchainIterator) Next() *Block {
	var block *Block
	_ = i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		encodedBlock := b.Get(i.currentHash)
		block = Deserialize(encodedBlock)
		return nil
	})
	i.currentHash = block.PrevHash
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
