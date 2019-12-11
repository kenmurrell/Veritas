# Veritas
 A basic blockchain written in Go

## Setup
```
$ go get -u github.com/boltdb/bolt
```

## Usage

### Creating and Using a Portfolio
A Portfolio stores multiple Wallets and is saved locally to *portfolio.dat* by default 
```
$ ./Veritas createportfolio
```
```
$ ./Veritas listaddresses
```
### Creating a Wallet
A Wallet stores your Public/Private key under an address (an encoded hash of your public key)
```
$ ./Veritas createwallet
```
```
$ ./Veritas balance -address 00964c39164c3936ab7aa1fd1a705b053c29af849166d5caf2
```
### Creating a Blockchain
A Blockchain requires an address to send the reward for mining the genesis block
```
$ ./Veritas createblockchain -address 00964c39164c3936ab7aa1fd1a705b053c29af849166d5caf2
```
### Sending Moolah
Sending money through a transaction requires your address, the recipient's address and an amount. The transaction is then stored in the mining queue and is not valid until it is mined
```
$ ./Veritas send -from 00964c39164c3936ab7aa1fd1a705b053c29af849166d5caf2 -to 004010e8a7e8a537cad8fc348baf7748625f77ab3f4206f302 -amount 1
```
### Mining Blocks
Mining a block requires an address to send the reward and the number of transactions to put into the block. The difficulty of the block increases the more transactions it contains
```
$ ./Veritas mine -address 00964c39164c3936ab7aa1fd1a705b053c29af849166d5caf2 -num 5
```
### Other Shenanigans
You can print the full blockchain using *printchain* 
```
$ ./Veritas printchain
```
