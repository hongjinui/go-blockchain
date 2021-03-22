package main

import (
	"fmt"
	"log"

	"github.com/hongjinui/go-blockchain/blockchain"
)

func main() {

	log.Println("main.go start!")
	log.Println("this is part1 ")

	bc := blockchain.NewBlockchain()

	bc.AddBlock("Send 1 BTC to hhh")
	bc.AddBlock("Send 2 more BTC to HHH")

	for _, block := range bc.GetBlocks() {
		fmt.Printf("Prev : hash : %x\n", block.PrevBlockHash)
		fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Println()
	}
}
