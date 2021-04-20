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
		fmt.Printf("=========== Block %x ===========\n", block.Hash)
		fmt.Printf("Prev. block : %x\n", block.PrevBlockHash)

		pow := b.NewProofOfWork(block)
		fmt.Printf("PoW : %s\n\n", strconv.FormatBool(pow.Validation()))

		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}

		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
