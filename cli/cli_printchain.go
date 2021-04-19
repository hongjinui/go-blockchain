package cli

import (
	"fmt"
	"strconv"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

func (cli *CLI) printChain() {
	bc := b.NewBlockchain("")
	defer bc.GetDB().Close()
	bci := bc.Iterator()

	for {
		block := bci.Next()
		fmt.Printf("Prev. hash : %x\n", block.PrevBlockHash)
		fmt.Printf("Hash : %x\n", block.Hash)

		pow := b.NewProofOfWork(block)
		fmt.Printf("PoW : %s\n", strconv.FormatBool(pow.Validation()))

		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}

		fmt.Println("")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
