package main

import (
	blockchain "blockchain/Blockchain"
	"fmt"
	"strconv"
)

func main() {
	chain := blockchain.InitBlockChain()
	chain.Addblock("First Block after genesis")
	chain.Addblock("Second Block after genesis")
	chain.Addblock("Third Block after genesis")
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash:%x\n \n", block.PrevHash)
		fmt.Printf("Data in Block:%s\n\n", block.Data)
		fmt.Printf("Hash:%x\n\n-\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
