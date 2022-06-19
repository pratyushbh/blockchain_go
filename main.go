package main

import (
	"os"

	"github.com/pratyushbh/blockchain_go/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}
