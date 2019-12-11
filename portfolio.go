package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Portfolio struct {
  Wallets map[string]*Wallet
}

func NewPortfolio() *Portfolio {
  p := Portfolio{}
  p.Wallets = make(map[string]*Wallet)
  return &p
}

func (p *Portfolio) CreateWallet() Wallet {
  wallet := NewWallet()
  address := fmt.Sprintf("%s", wallet.GetAddress())
  p.Wallets[address] = wallet
	return *wallet
}

func (p Portfolio) GetWallet(address string) Wallet {
	return *p.Wallets[address]
}

func (p *Portfolio) GetAllAddresses() []string {
	var addresses []string
	for address := range p.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func LoadPortfolio(file string) (*Portfolio, error) {
  portfolio := Portfolio{}
  portfolio.Wallets = make(map[string]*Wallet)
  if _, err := os.Stat(file); os.IsNotExist(err) {
    log.Panic(err)
  }
  content, _ := ioutil.ReadFile(file)
  gob.Register(elliptic.P256())
  decoder := gob.NewDecoder(bytes.NewReader(content))
  var temp Portfolio
  err := decoder.Decode(&temp)
  portfolio.Wallets = temp.Wallets
  return &portfolio, err
}

func (p Portfolio) SaveToFile(file string) {
  var content bytes.Buffer
  gob.Register(elliptic.P256())
  encoder := gob.NewEncoder(&content)
  _ = encoder.Encode(&p)
  err := ioutil.WriteFile(file, content.Bytes(), 0644)
  if(err != nil) {
    log.Panic(err)
  }

}
