package cli

import (
	"fmt"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) send(from, to string, amount int) {
	bc := b.NewBlockchain(from)
	defer bc.GetDB().Close()

	tx := b.NewUTXOTransaction(from, to, amount, bc)
	bc.MindBlock([]*b.Transaction{tx})

	fmt.Println("Success")
}
