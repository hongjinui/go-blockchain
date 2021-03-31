package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	b "github.com/hongjinui/go-blockchain/blockchain"
)

type CLI struct {
	bc *b.Blockchain
}

func (cli *CLI) SetBC(bc *b.Blockchain) {
	cli.bc = bc
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage :")
	fmt.Println("	getBalance -address ADDRESS -Get balance of ADDRESS")
	fmt.Println("	createblockchin -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("	printchain - print all the blocks of the blockchain")

	fmt.Println("	send -from FROM -to TO -amount AMOUNT -send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchintCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchintCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := flag.String("from", "", "Source wallat address")
	sendTo := flag.String("to", "", "Destination wallat address")
	sendAmount := flag.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchintCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}
	if createBlockchintCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchintCmd.Usage()
			os.Exit(1)
		}
		cli.CreateBlockchain(*createBlockchainAddress)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

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
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
func (cli *CLI) send(from, to string, amount int) {
	bc := b.NewBlockchain(from)
	defer bc.GetDB().Close()

	tx := b.NewUTOXTransaction(from, to, amount, bc)
	cli.bc.AddBlock([]*b.Transaction{tx})

	fmt.Println("Success")
}

func (cli *CLI) getBalance(address string) {
	bc := b.NewBlockchain(address)
	defer bc.GetDB().Close()

	balance := 0
	utxs := cli.bc.FindUnspentTransaction(address)

	for _, tx := range utxs {
		for _, out := range tx.Vout {
			if out.Unlock(address) {
				balance += out.Value
			}
		}
	}
	fmt.Printf("Balnce of '%s' : %d\n")

}
func (cli *CLI) CreateBlockchain(address string) {
	bc := b.CreateBlockchain(address)
	bc.GetDB().Close()
	fmt.Println("DONE!!")

}
