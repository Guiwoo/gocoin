package main

import (
	"fmt"

	"github.com/guiwoo/gocoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockChain()
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")
	chain.AddBlock("Fourth Block")
	for _, block := range chain.AllBlocks() {
		fmt.Printf("\nData: %s\n", block.Data)
		fmt.Printf("Data: %s\n", block.Hash)
		fmt.Printf("Data: %s\n", block.PrevHash)
	}
}
