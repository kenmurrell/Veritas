package main

import (
  "bytes"
  "time"
  "encoding/gob"
  "crypto/sha256"
)

type Block struct {
  Timestamp   int64
  Transactions []*Transaction
  PrevHash    []byte
  Hash        []byte
  Nonce       int
}

func NewGenesisBlock(baseTransaction *Transaction) *Block {
  return NewBlock([]*Transaction{baseTransaction}, []byte{})
}

func NewBlock(transaction []*Transaction, prevHash []byte) *Block {
  timestamp := time.Now().Unix()
  block := &Block{timestamp, transaction, prevHash, []byte{}, 0}
  pow := NewProofOfWork(block)
  nonce, hash := pow.Generate()
  block.Nonce = nonce
  block.Hash = hash[:]
  return block
}

func (b *Block) HashBlock(nonce int, target int) []byte {
  var txIDs [][]byte
  for _,tx := range b.Transactions{
    txIDs = append(txIDs, tx.ID)
  }
  var txHash [32]byte
  txHash = sha256.Sum256(bytes.Join(txIDs,[]byte{}))
  blockData := [][]byte{b.PrevHash, txHash[:], IntToHex(b.Timestamp), IntToHex(int64(target)), IntToHex(int64(nonce)),}
  data := bytes.Join(blockData, []byte{},)
  return data
}

func (b *Block) Serialize() []byte {
  var result bytes.Buffer
  encoder := gob.NewEncoder(&result)
  _ = encoder.Encode(b)
  return result.Bytes()
}

func Deserialize(d []byte) *Block {
  var block Block
  decoder := gob.NewDecoder(bytes.NewReader(d))
  _ = decoder.Decode(&block)
  return &block
}
