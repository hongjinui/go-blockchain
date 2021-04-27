package cli

import (
	"fmt"

	s "github.com/hongjinui/go-blockchain/server"
)

func (cli *CLI) startNode(nodeID string) {
	fmt.Printf("Starting node %s\n", nodeID)
	s.StartServer(nodeID)
}
