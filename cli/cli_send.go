package cli

import (
	"fmt"
	"log"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) send(from, to string, amount int) {
	if !b.ValidateAddress(from) {
		log.Panic("ERROR : Sender address is not valid")
	}

	if !b.ValidateAddress(to) {
		log.Panic("ERROR : Recipient address is not valid")
	}
	bc := b.NewBlockchain(from)
	defer bc.GetDB().Close()

	tx := b.NewUTXOTransaction(from, to, amount, bc)
	bc.MindBlock([]*b.Transaction{tx})

	fmt.Println("Success")
}
