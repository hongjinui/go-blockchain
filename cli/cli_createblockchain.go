package cli

import (
	"fmt"
	"log"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) createBlockchain(address string) {
	if !b.ValidateAddress(address) {
		log.Panic("ERROR : Address is not valid")
	}

	bc := b.CreateBlockchain(address)
	bc.GetDB().Close()
	fmt.Println("Done!")
}
