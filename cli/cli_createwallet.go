package cli

import (
	"fmt"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) createWallet() {
	wallets, _ := b.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address : %s\n", address)
}
