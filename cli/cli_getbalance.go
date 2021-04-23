package cli

import (
	"fmt"
	"log"

	b "github.com/hongjinui/go-blockchain/blockchain"
	"github.com/hongjinui/go-blockchain/utils"
)

func (cli *CLI) getBalance(address string) {
	if !b.ValidateAddress(address) {
		log.Panic("ERROR : Address is not valid")
	}

	bc := b.NewBlockchain()
	UTXOSet := b.UTXOSet{bc}

	defer bc.GetDB().Close()

	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balnce of '%s' : %d\n", address, balance)

}
