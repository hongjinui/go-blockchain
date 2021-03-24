package main

import (
	b "github.com/hongjinui/go-blockchain/blockchain"
	c "github.com/hongjinui/go-blockchain/cli"
)

func main() {

	bc := b.NewBlockchain()
	db := bc.GetDB()
	defer db.Close()
	cli := c.CLI{}
	cli.SetBC(bc)
	cli.Run()

}
