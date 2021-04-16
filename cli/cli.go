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
}

func (cli *CLI) createBlockchain(address string) {
	bc := b.CreateBlockchain(address)
	bc.GetDB().Close()
	fmt.Println("Done!")
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
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
	createblockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createblockchainAddress := createblockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallat address")
	sendTo := sendCmd.String("to", "", "Destination wallat address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createblockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressCmd.Parse(os.Args[2:])
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
	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	if listAddressCmd.Parsed() {
		cli.listAddresses()
	}
	if createblockchainCmd.Parsed() {
		if *createblockchainAddress == "" {
			createblockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createblockchainAddress)
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
		// fmt.Println("Transactions : [")
		// fmt.Println("]")
		fmt.Println("")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
func (cli *CLI) send(from, to string, amount int) {
	bc := b.NewBlockchain(from)
	defer bc.GetDB().Close()

	tx := b.NewUTXOTransaction(from, to, amount, bc)
	bc.MindBlock([]*b.Transaction{tx})

	fmt.Println("Success")
}

func (cli *CLI) getBalance(address string) {
	bc := b.NewBlockchain(address)
	defer bc.GetDB().Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balnce of '%s' : %d\n", address, balance)

}

func (cli *CLI) createWallet() {
	wallets := b.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address : %s\n", address)
}

func (cli *CLI) listAddresses() {
	wallets := b.NewWallets()
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}

}
