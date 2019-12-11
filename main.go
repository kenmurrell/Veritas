package main

import (
	"fmt"
)

func main() {

	p1, p2 := MyWallets()
	bc := MyBlockchain(p1)
	balance1 := bc.CalculateBalance(p1.GetAddress())
	fmt.Printf("-->p1: %d\n", balance1)
	balance2 := bc.CalculateBalance(p2.GetAddress())
	fmt.Printf("-->p2: %d\n", balance2)

	Send(bc, p1, p2.GetAddress(), 1)
	Send(bc, p1, p2.GetAddress(), 2)
	Mine(bc, p2.GetAddress(), 3)
	balance1 = bc.CalculateBalance(p1.GetAddress())
	fmt.Printf("-->p1: %d\n", balance1)
	balance2 = bc.CalculateBalance(p2.GetAddress())
	fmt.Printf("-->p2: %d\n", balance2)
	bc.Printchain()
}
