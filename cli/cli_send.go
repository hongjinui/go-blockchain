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
	bc := b.NewBlockchain()
	UTXOSet := b.UTXOSet{bc}

	defer bc.GetDB().Close()

	tx := b.NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := b.NewCoinbaseTX(from, "")
	txs := []*b.Transaction{cbTx, tx}
	bc.MindBlock(txs)

	bc.MindBlock([]*b.Transaction{tx})

	fmt.Println("Success")
}
