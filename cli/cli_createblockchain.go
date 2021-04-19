package cli

import (
	"fmt"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) createBlockchain(address string) {
	bc := b.CreateBlockchain(address)
	bc.GetDB().Close()
	fmt.Println("Done!")
}
