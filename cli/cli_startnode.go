package cli

import (
	"fmt"
	"log"

	b "github.com/hongjinui/go-blockchain/blockchain"
	s "github.com/hongjinui/go-blockchain/server"
)

func (cli *CLI) startNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if b.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}

	s.StartServer(nodeID, minerAddress)
}
