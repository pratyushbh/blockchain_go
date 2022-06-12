package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/tensor-programming/golang-blockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.Blockchain
}

func (cli *CommandLine) Instructions() {
	fmt.Println("Usage:")
	fmt.Println("add -block BLOCK_DATA -add a block to the Chain")
	fmt.Println("print - to print the blocks in the chain")
}
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.Instructions()
		runtime.Goexit()
	}
}
func (cli *CommandLine) addblock(data string) {
	cli.blockchain.Addblock(data)
	fmt.Println("Added Block!")
}
func (cli *CommandLine) PrintBlockchain() {
	iter := cli.blockchain.Iterator()
	for {
		block := iter.Next()
		fmt.Printf("Previous Hash:%x\n \n", block.PrevHash)
		fmt.Printf("Data in Block:%s\n\n", block.Data)
		fmt.Printf("Hash:%x\n\n-\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevHash) == 0 {
			break
		}
	}
}
func (cli *CommandLine) run() {
	cli.validateArgs()
	addBlockArgs := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockArgs.String("block", "", "Block Data")
	switch os.Args[1] {
	case "add":
		err := addBlockArgs.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.Instructions()
		runtime.Goexit()
	}
	if addBlockArgs.Parsed() {
		if *addBlockData == "" {
			cli.Instructions()
			runtime.Goexit()
		}
		cli.addblock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		cli.PrintBlockchain()
	}
}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()
	cli := CommandLine{chain}
	cli.run()
}
