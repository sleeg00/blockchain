package main

import (
	"fmt"
	"log"
	"strconv"
)

func (cli *CLI) printChain(nodeID string) {
	log.Println("printChain")
	bc := NewBlockchain(nodeID)
	log.Println("??")
	defer bc.db.Close()

	bci := bc.Iterator()

	for {

		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		UTXOSet := UTXOSet{Blockchain: bc}
		fmt.Println(UTXOSet.CountTransactions())
		fmt.Println(UTXOSet.Blockchain.tip)
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
