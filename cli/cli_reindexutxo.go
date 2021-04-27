package cli

import (
	"fmt"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) reindexUTXO(nodeID string) {
	bc := b.NewBlockchain(nodeID)
	UTXOSet := b.UTXOSet{bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)

}
