package main

import (
	"flag"
	"fmt"
	"log"
	"os"
  "gopkg.in/yaml.v2"
  "io/ioutil"
)

type CLI struct {
  conf *Config
}

type Config struct {
  BlockchainFile string `yaml:"blockchainFile"`
  PortfolioFile string `yaml:"portfolioFile"`
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func start() *CLI {
  return &CLI{loadConfig()}
}

func (cli *CLI) Run() {
	cli.validateArgs()
  //Balance command
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
  bAddr := balanceCmd.String("address", "", "Address of the account")
  //Send Command
  sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
  fromAddr := sendCmd.String("from", "", "Address to send from")
	toAddr := sendCmd.String("to", "", "Address to send to")
	amount := sendCmd.Int("amount", 0, "Amount to send")
  //Mine Command
  mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)
  mineAddr := mineCmd.String("address", "", "Address to reward for mining")
  numTx := mineCmd.Int("num", 1, "Number of transactions to put into block to mine")
  //Create Bc Command
	createBcCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
  createAddr := createBcCmd.String("address", "", "Address to reward for genesis block")
  //Create Portfolio Command
  createPCmd := flag.NewFlagSet("createportfolio", flag.ExitOnError)
  //Create Wallet Command
	createWCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	//List Addresses Command
  listAddrCmd:= flag.NewFlagSet("listaddresses", flag.ExitOnError)
  //Print blockchain Command
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	switch os.Args[1] {
	case "balance":
		_ = balanceCmd.Parse(os.Args[2:])
  case "send":
		_ = sendCmd.Parse(os.Args[2:])
  case "mine":
    _ = mineCmd.Parse(os.Args[2:])
	case "createblockchain":
		_ = createBcCmd.Parse(os.Args[2:])
  case "createportfolio":
    _ = createPCmd.Parse(os.Args[2:])
	case "createwallet":
		_ = createWCmd.Parse(os.Args[2:])
	case "listaddresses":
    _ = listAddrCmd.Parse(os.Args[2:])
	case "printchain":
		_ = printChainCmd.Parse(os.Args[2:])
	default:
		fmt.Printf("--> Unrecognized command")
    cli.PrintUsage()
		os.Exit(1)
	}
	if balanceCmd.Parsed() {
		if *bAddr == "" {
			balanceCmd.Usage()
			os.Exit(1)
		}
		cli.GetBalance(*bAddr)
	}
  if sendCmd.Parsed() {
		if *toAddr == "" || *fromAddr == "" || *amount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.Send(*toAddr, *fromAddr, *amount)
	}
  if mineCmd.Parsed() {
    if *mineAddr == "" {
      mineCmd.Usage()
      os.Exit(1)
    }
    cli.Mine(*mineAddr, *numTx)
  }
	if createBcCmd.Parsed() {
		if *createAddr == "" {
			createBcCmd.Usage()
			os.Exit(1)
		}
		cli.CreateBlockchain(*createAddr)
	}
  if createPCmd.Parsed() {
    cli.CreatePortfolio()
  }
	if createWCmd.Parsed() {
		cli.CreateWallet()
	}
	if listAddrCmd.Parsed() {
		cli.ListAddresses()
	}
	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}

func (cli *CLI) GetBalance(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address not valid!\n")
	}
	bc, err := LoadBlockchain(cli.conf.BlockchainFile)
  if err != nil {
    log.Panic(err.Error())
  }
  defer bc.Close()
  balance := bc.CalculateBalance(address)
	fmt.Printf("--> balance of '%s': %d\n", address, balance)
}

func (cli *CLI) Send(fromAddr, toAddr string, amount int) {
	if !ValidateAddress(fromAddr) {
		log.Panic("ERROR: Sender is not valid\n")
	}
	if !ValidateAddress(toAddr) {
		log.Panic("ERROR: Recipient is not valid\n")
	}
	bc, err := LoadBlockchain(cli.conf.BlockchainFile)
  if err != nil {
    log.Panic(err.Error())
  }
  defer bc.Close()
  portfolio, err := LoadPortfolio(cli.conf.PortfolioFile)
  if err != nil {
    log.Panic(err.Error())
  }
  wallet := portfolio.GetWallet(fromAddr)
  tx := NewTransaction(wallet.GetAddress(), toAddr, amount)
	err = tx.Sign(&wallet)
  if err != nil {
    log.Panic(err.Error())
  }
  if tx.Validate() {
    bc.AddToQueue(tx)
  } else {
    log.Panic("ERROR: Cannot validate transaction")
  }
	fmt.Println("--> complete\n")
}

func (cli *CLI) Mine(address string, num int) {
  if !ValidateAddress(address) {
		log.Panic("ERROR: Sender is not valid\n")
	}
  bc, err := LoadBlockchain(cli.conf.BlockchainFile)
  if err != nil {
    log.Panic(err.Error())
  }
  defer bc.Close()
  err = bc.MineBlock(num)
  if err != nil {
    log.Panic(err.Error())
  }
  fmt.Printf("--> block mined\n")
}

func (cli *CLI) CreateBlockchain(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid\n")
	}
	bc, err := CreateBlockchain(cli.conf.BlockchainFile, address)
  if err != nil {
    log.Panic(err.Error())
  }
  defer bc.Close()
	fmt.Println("--> blockchain created\n")
}

func (cli *CLI) CreatePortfolio() {
  portfolio := NewPortfolio()
  portfolio.SaveToFile("portfolio.dat")
  fmt.Printf("--> portfolio created\n")
}

func (cli *CLI) CreateWallet() {
  portfolio, err := LoadPortfolio(cli.conf.PortfolioFile)
  if err != nil {
    log.Panic(err.Error())
  }
	wallet := portfolio.CreateWallet()
	err = portfolio.SaveToFile("portfolio.dat")
  if err != nil {
    log.Panic(err.Error())
  }
	fmt.Printf("--> new address: %s\n", wallet.GetAddress())
}

func (cli *CLI) ListAddresses() {
	portfolio, err := LoadPortfolio(cli.conf.PortfolioFile)
  if err != nil {
    log.Panic(err.Error())
  }
  addresses := portfolio.GetAllAddresses()
  fmt.Printf("--> all addresses:\n")
	for _, address := range addresses {
		fmt.Printf("==> %s\n",address)
	}
}

func (cli *CLI) PrintChain() {
	bc, err := LoadBlockchain(cli.conf.BlockchainFile)
  if err != nil {
    log.Panic(err.Error())
  }
	defer bc.Close()
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

func (cli* CLI) PrintUsage(){
  fmt.Printf("--> see README.md for usage\n")
}

func loadConfig() *Config {
  var config Config
  src, err := ioutil.ReadFile("config.yaml")
  if err != nil {
    log.Panic(err)
  }
  err = yaml.Unmarshal(src, &config)
  if err != nil {
    log.Panic(err.Error())
  }
  return &config
}
