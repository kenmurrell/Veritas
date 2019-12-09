package main

import (
  "bytes"
  "crypto/sha256"
  "fmt"
  "log"
  "math"
  "math/big"
  "encoding/binary"
)

type ProofOfWork struct {
  block *Block
  target *big.Int
}

const targetBits = 20

func NewProofOfWork(b *Block) *ProofOfWork {
  target := big.NewInt(1)
  target.Lsh(target, uint(256 - targetBits))
  pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) Generate() (int, []byte) {
  var hashInt big.Int
  var hash [32]byte
  nonce := 0
  fmt.Printf("--> mining a new block...")
  for nonce < math.MaxInt64 {
    data := pow.block.HashBlock(nonce, targetBits)
    hash = sha256.Sum256(data)
    fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
  }
  fmt.Printf("\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.block.HashBlock(pow.block.Nonce, targetBits)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
