package main

import (
	"strconv"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"crypto/sha256"
	"errors"
	"encoding/gob"
	"bytes"
	"fmt"
	mathrand "math/rand"
	"time"
)

type Transaction struct {
	ID   []byte
	Timestamp   int64
	FromAddr string
	ToAddr string
	Amount int
	Signature []byte
	PublicKey []byte
}

const subsidy = 10

// TODO: Convert this to UXTO when i get the chance?
func NewGenesisTransaction(toAddr, startData string) *Transaction {
	id := []byte("GENESIS")
	timestamp := time.Now().UnixNano()
	tx := Transaction{id, timestamp, startData, toAddr, subsidy, nil, nil}
	return &tx
}

func NewTransaction(fromAddr, toAddr string, amount int) *Transaction {
	if !ValidateAddress(fromAddr) && !ValidateAddress(toAddr) {
		panic(errors.New("Addresses not valid"))
	}
	timestamp := time.Now().UnixNano()
	mathrand.Seed(timestamp)
	hash := sha256.Sum256([]byte(fromAddr + toAddr + strconv.Itoa(mathrand.Intn(1000))))
	tx := &Transaction{hash[:], timestamp, fromAddr, toAddr, amount, nil, nil}
	return tx
}

func (tx *Transaction) Sign(wallet *Wallet) error {
	if wallet.GetAddress() != tx.FromAddr {
		return errors.New("ERROR: criminal detected")
	}
	r, s, err := ecdsa.Sign(rand.Reader, &wallet.PrivateKey, tx.ID)
	if err != nil {
		return errors.New("ERROR: criminal detected")
	}
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = signature
	tx.PublicKey = wallet.PublicKey
	return nil
}

func (tx *Transaction) ValidateSignature() bool {
		//regenerate r, s
		r := big.Int{}
		s := big.Int{}
		sigLen := len(tx.Signature)
		r.SetBytes(tx.Signature[:(sigLen/2)])
		s.SetBytes(tx.Signature[(sigLen/2):])
		//regenerate public key
		curve := elliptic.P256()
		x := big.Int{}
		y := big.Int{}
		keyLen := len(tx.PublicKey)
		x.SetBytes(tx.PublicKey[:(keyLen/2)])
		y.SetBytes(tx.PublicKey[(keyLen/2):])
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		status := ecdsa.Verify(&rawPubKey, tx.ID, &r, &s)
		return status
}

func (tx *Transaction) Serialize() []byte {
  var result bytes.Buffer
  encoder := gob.NewEncoder(&result)
  _ = encoder.Encode(tx)
  return result.Bytes()
}

func DeserializeToTX(enc []byte) *Transaction {
  var tx Transaction
  decoder := gob.NewDecoder(bytes.NewReader(enc))
  _ = decoder.Decode(&tx)
  return &tx
}

func (tx *Transaction) ToString() {
	fmt.Printf("\t===== Transaction %x =====\n", tx.ID)
	fmt.Printf("\tTimestamp: %s\n", time.Unix(0, tx.Timestamp))
	fmt.Printf("\tFrom Address: %s\n", tx.FromAddr)
	fmt.Printf("\tTo Address: %s\n", tx.ToAddr)
	fmt.Printf("\tAmount: %d\n", tx.Amount)
}
