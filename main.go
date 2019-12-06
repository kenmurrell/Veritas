package main

import (
	"fmt"
)

func main() {
	// portfolio := NewPortfolio()
	// addr1 := portfolio.CreateWallet()
	// person1 := portfolio.GetWallet(addr1)
	// addr2 := portfolio.CreateWallet()
	// person2 := portfolio.GetWallet(addr2)
	// portfolio.SaveToFile("port.dat")

	portfolio, _ := LoadPortfolio("port.dat")
	person1 := portfolio.GetWallet("001b950fbf3d818d9b740907b8c290178548cc48c275b49fd3")
	person2 := portfolio.GetWallet("0003ee92244c20027b9406f46fcff52297900aa9528a2675ac")

	fmt.Printf("Person1 wallet address: %s %t\n", person1.GetAddress(), ValidateAddress(person1.GetAddress()))
	fmt.Printf("Person2 wallet address: %s %t\n", person2.GetAddress(), ValidateAddress(person2.GetAddress()))
	fmt.Printf("Fakeper wallet address: %s %t\n", "0003ee92244c20027b9406f46fcff52297900aa9528a2675ab", ValidateAddress("0003ee92244c20027b9406f46fcff52297900aa9528a2675ab"))

	 //_ = CreateBlockchain(person1.GetAddress())

	 _ = LoadBlockchain(person1.GetAddress())

	 

}
