package main

import (
	"os"

	blockchain "github.com/pratyushbh/blockchain_go/wallet"
)

func main() {
	defer os.Exit(0)
	// cmd := cli.CommandLine{}
	// cmd.Run()
	w := blockchain.MakeWallet()
	w.Address()
}
