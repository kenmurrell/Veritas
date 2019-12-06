package main

type Transaction struct {
	ID   []byte
	fromAddr string
	toAddr string
	amount int
}

//TODO: SIGN TRANSACTIONS!!!

const subsidy = 10

// TODO: Convert this to UXTO when i get the chance?
func NewGenesisTransaction(toAddr, startData string) *Transaction {
	tx := Transaction{nil, startData, toAddr, subsidy}
	return &tx
}
