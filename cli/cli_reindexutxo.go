package cli

import (
	"fmt"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) reindexUTXO() {
	bc := b.NewBlockchain()
	UTXOSet := b.UTXOSet{bc}
	UTXOSet.Reindex()

	count := UTXOSet.GetCount()
	fmt.Printf("Done! There are %d transactions in the UTXO set.", count)

}
