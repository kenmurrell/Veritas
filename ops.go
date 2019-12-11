package main

import (
  "fmt"
)


func MyWallets() (*Wallet, *Wallet) {
  portfolio, _ := LoadPortfolio("port.dat")
	person1 := portfolio.GetWallet("001b98d6a7e8a226884fc9551ef63a96fd21df7d8a8d3bb772")
	person2 := portfolio.GetWallet("009b0163effe109fd1b707a0c092871de23f9ae44b45503ff8")
  return &person1, &person2
}

func CreateMyPorfolio() (string, string) {
  portfolio := NewPortfolio()
	addr1 := portfolio.CreateWallet()
	addr2 := portfolio.CreateWallet()
	portfolio.SaveToFile("port.dat")
  return addr1, addr2
}

func MyBlockchain(w *Wallet) *Blockchain {
  bc := CreateBlockchain(w.GetAddress())
  return bc
}

func Send(bc *Blockchain, w *Wallet, addr string, amount int) {
  tx := NewTransaction(w.GetAddress(), addr, amount)
	_ = tx.Sign(w)
  if tx.Validate() {
    bc.AddToQueue(tx)
  }
}

func Mine(bc *Blockchain, addr string, num int) {
  bc.MineBlock(num)
}

func (bc *Blockchain) Printchain() {
  fmt.Printf("--> Printing chain...\n")
  bci := bc.Iterator()
  for {
    block := bci.Next()
    block.ToString()
    for _, tx := range block.Transactions {
			tx.ToString()
		}
    if len(block.PrevHash) == 0 {
      break
    }
  }
}
