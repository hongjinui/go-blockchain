package cli

import (
	"fmt"

	b "github.com/hongjinui/go-blockchain/blockchain"
	"github.com/hongjinui/go-blockchain/utils"
)

func (cli *CLI) getBalance(address string) {
	bc := b.NewBlockchain(address)
	defer bc.GetDB().Close()

	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balnce of '%s' : %d\n", address, balance)

}
