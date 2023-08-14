package main

import (
	"context"
	"fmt"
	"log"

	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

func RsEncoding(count int32, f int32, NF int32) {

	for i := 0; i < len(knownNodes); i++ {
		serverAddress := fmt.Sprintf("localhost:%s", knownNodes[i][10:])
		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()

		client := blockchain.NewBlockchainServiceClient(conn)
		cli := CLI{
			nodeID:     knownNodes[i][10:],
			blockchain: client,
		}

		cli.blockchain.RSEncoding(context.Background(), &blockchain.RSEncodingRequest{
			NodeId: knownNodes[i][10:],
			Count:  count - 1,
			F:      f,
			NF:     NF,
		})
	}

}
