package main

import (
	"context"
	"fmt"
	"log"

	blockchain "github.com/sleeg00/blockchain/proto"
)

/*
	func (cli *CLI) createBlockchain(address, nodeID string) {
		if !ValidateAddress(address) {
			log.Panic("ERROR: Address is not valid")
		}
		bc := CreateBlockchain(address, nodeID)
		defer bc.db.Close()

		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()

		fmt.Println("Done!")
	}
*/
func (cli *CLI) createBlockchain(address, nodeID string) {
	request := &blockchain.CreateBlockchainRequest{
		Address: address,
		NodeId:  nodeID,
	}

	response, err := cli.blockchain.CreateBlockchain(context.Background(), request)
	if err != nil {
		log.Fatalf("Error calling CreateBlockchain RPC : %v", err)
	}

	if response.Response == "Success" {
		fmt.Println("Success!")
	} else {
		fmt.Println("failed.")
	}
}
