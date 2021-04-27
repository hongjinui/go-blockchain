package cli

import (
	"fmt"
	"log"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) send(from, to string, amount int, nodeID string) {
	if !b.ValidateAddress(from) {
		log.Panic("ERROR : Sender address is not valid")
	}

	if !b.ValidateAddress(to) {
		log.Panic("ERROR : Recipient address is not valid")
	}
	bc := b.NewBlockchain(nodeID)
	UTXOSet := b.UTXOSet{bc}

	defer bc.GetDB().Close()

	wallets, err := b.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	tx := b.NewUTXOTransaction(&wallet, to, amount, &UTXOSet)
	cbTx := b.NewCoinbaseTX(from, "")
	txs := []*b.Transaction{cbTx, tx}

	newBlock := bc.MindBlock(txs)
	UTXOSet.Update(newBlock)

	bc.MindBlock([]*b.Transaction{tx})

	fmt.Println("Success")
}
