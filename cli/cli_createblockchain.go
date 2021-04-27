package cli

import (
	"fmt"
	"log"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) createBlockchain(address, nodeID string) {
	if !b.ValidateAddress(address) {
		log.Panic("ERROR : Address is not valid")
	}

	bc := b.CreateBlockchain(address, nodeID)

	defer bc.GetDB().Close()

	UTXOSet := b.UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}
