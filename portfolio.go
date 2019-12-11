package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"errors"
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
    return nil, errors.New("ERROR: Portfolio file does not exist\n")
  }
  content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.New("ERROR: Cannot read portfolio file\n")
	}
  gob.Register(elliptic.P256())
  decoder := gob.NewDecoder(bytes.NewReader(content))
  var temp Portfolio
  err = decoder.Decode(&temp)
	if err != nil {
		return nil, errors.New("ERROR: Cannot decode portfolio file\n")
	}
  portfolio.Wallets = temp.Wallets
  return &portfolio, nil
}

func (p Portfolio) SaveToFile(file string) error {
  var content bytes.Buffer
  gob.Register(elliptic.P256())
  encoder := gob.NewEncoder(&content)
  err := encoder.Encode(&p)
	if err != nil {
		return errors.New("ERROR: Cannot encode portfolio\n")
	}
  err = ioutil.WriteFile(file, content.Bytes(), 0644)
	if err != nil {
		return errors.New("ERROR: Cannot save portfolio to filename\n")
	}
  return nil

}
