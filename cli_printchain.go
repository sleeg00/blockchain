package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

func (cli *CLI) printChain(nodeID string) {

	bc := NewBlockchainRead(nodeID)

	defer bc.db.Close()

	bci := bc.Iterator()

	data := make([][]byte, 10)

	var Height int

	for x := 0; ; x++ {

		block, err := bci.Next()
		if x == 0 {
			Height = block.Height
		}
		checkErr(err)
		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 || x == Height%7 {
			log.Println("xxxxxxxx", x)
			break
		}
	}

	count := Height / 7

	blockSize := 5120 // 고정된 샤드 크기
	for i := 0; i < Height/7; i++ {
		count--
		for k := 0; k < len(knownNodes); k++ {
			data[k] = make([]byte, blockSize)

			serverAddress := fmt.Sprintf("localhost:%s", knownNodes[k][10:])

			conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Failed to connect to gRPC server: %v", err)
			}
			defer conn.Close()

			client := blockchain.NewBlockchainServiceClient(conn)
			cli := CLI{
				nodeID:     knownNodes[k][10:],
				blockchain: client,
			}
			request := &blockchain.GetShardRequest{
				NodeId: knownNodes[k][10:],
				Height: int32(count),
			}

			response, err := cli.blockchain.GetShard(context.Background(), request)

			bytes := response.Bytes
			//log.Println(DeserializeBlock(bytes))
			data[k] = bytes

		}
		for x := 0; x < 7; x++ {

			var result []byte

			for y := 0; y < len(data[x]); y++ {

				result = append(result, data[x][y])

			}

			block := DeserializeBlock(result)
			fmt.Printf("============ Block %x ============\n", block.Hash)
			fmt.Printf("Height: %d\n", block.Height)
			fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
			pow := NewProofOfWork(block)
			fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
			for _, tx := range block.Transactions {
				fmt.Println(tx)
			}
			fmt.Printf("\n\n")

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}
	}
}
