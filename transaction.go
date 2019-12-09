package main

import (
	"strconv"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"crypto/sha256"
	"errors"
)

type Transaction struct {
	ID   []byte
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
	tx := Transaction{id, startData, toAddr, subsidy, nil, nil}
	return &tx
}

func NewTransaction(fromAddr, toAddr string, amount int) *Transaction {
	//verify addreses here
	hash := sha256.Sum256([]byte(fromAddr + toAddr + strconv.Itoa(amount)))
	tx := &Transaction{hash[:], fromAddr, toAddr, amount, nil, nil}
	return tx
}

func (tx *Transaction) Sign(wallet *Wallet) error {
	if wallet.GetAddress() != tx.FromAddr {
		return errors.New("wrong")
	}
	r, s, err := ecdsa.Sign(rand.Reader, &wallet.PrivateKey, tx.ID)
	if err != nil {
		return err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = signature
	tx.PublicKey = wallet.PublicKey
	return nil
}

func (tx *Transaction) Validate() bool {
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
