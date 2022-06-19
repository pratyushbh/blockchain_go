package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/pratyushbh/blockchain_go"
)

type CommandLine struct{}

func (cli *CommandLine) Instructions() {
	fmt.Println("Usage:")
	fmt.Println("getbalance -address ADDRESS- get balance for the account")
	fmt.Println("createblockchain -address ADDRESS creates a blockchain")
	fmt.Println("printchain - Prints the blocks in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT - Send amount from one account to another")
}
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.Instructions()
		runtime.Goexit()
	}
}
func (cli *CommandLine) PrintBlockchain() {
	chain := blockchain_go.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()
	for {
		block := iter.Next()
		fmt.Printf("\nPrevious Hash:%x\n \n", block.PrevHash)
		fmt.Printf("Hash:%x\n\n-\n", block.Hash)
		pow := blockchain_go.NewProof(block)
		fmt.Printf("PoW:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevHash) == 0 {
			break
		}
	}
}
func (cli *CommandLine) getbalance(address string) {
	chain := blockchain_go.ContinueBlockChain(address)
	defer chain.Database.Close()
	balance := 0
	UTXOs := chain.FindUTXO(address)
	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of %s:%d\n", address, balance)
}
func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain_go.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain_go.ContinueBlockChain(from)
	defer chain.Database.Close()
	tx := blockchain_go.NewTransaction(from, to, amount, chain)
	chain.Addblock([]*blockchain_go.Transaction{tx})
	fmt.Println("Success!")
}
func (cli *CommandLine) run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
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
		cli.Instructions()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getbalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.PrintBlockchain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
