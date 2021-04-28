package cli

import (
	"fmt"
	"log"

	b "github.com/hongjinui/go-blockchain/blockchain"
	s "github.com/hongjinui/go-blockchain/server"
)

func (cli *CLI) send(from, to string, amount int, nodeID string, mineNow bool) {
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

	if mineNow {
		txs := []*b.Transaction{cbTx, tx}
		newBlock := bc.MindBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		knownNodes := s.GetKnownNodes()
		s.SendTx(knownNodes[0], tx)
		s.SendTx(knownNodes[0], cbTx)
	}

	// newBlock := bc.MindBlock(txs)
	// UTXOSet.Update(newBlock)

	fmt.Println("Success")
}
