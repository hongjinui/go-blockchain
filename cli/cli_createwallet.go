package cli

import (
	"fmt"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := b.NewWallets(nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)

	fmt.Printf("Your new address : %s\n", address)
}
