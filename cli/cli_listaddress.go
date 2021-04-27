package cli

import (
	"fmt"
	"log"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) listAddresses(nodeID string) {
	wallets, err := b.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddresses()
	for _, address := range addresses {
		fmt.Println(address)
	}

}
