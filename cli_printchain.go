package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/klauspost/reedsolomon"
	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

func (cli *CLI) printChain(nodeID string) {
	nodeId, err := strconv.Atoi(nodeID)
	bc := NewBlockchain(nodeID)

	defer bc.db.Close()

	bci := bc.Iterator()
	block, err := bci.Next()
	if err != nil {
		log.Println(err)
	}
	for i := 0; i <= block.Height/7+(block.Height%7); i++ {

		data := make([][]byte, 10)
		if i < block.Height/7 { //블럭을 가져오고 만약에 죽었다면 그걸 바탕으로 디코딩해서 원본 파일을 생성
			log.Println("청크 찾기 ")
			count := block.Height / 7

			blockSize := 5120 // 고정된 샤드 크기
			enc, err := reedsolomon.New(7, 3)
			checkErr(err)
			log.Println("1")
			for i := 0; i < len(knownNodes); i++ {
				log.Println("2")
				data[i] = make([]byte, blockSize)

				if knownNodes[i][10:] != nodeID {
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
					request := &blockchain.GetShardRequest{
						NodeId: knownNodes[i][10:],
						Height: int32(count - 1),
					}
					log.Println("3")
					response, err := cli.blockchain.GetShard(context.Background(), request)
					log.Println("4")
					bytes := response.Bytes
					//log.Println(DeserializeBlock(bytes))
					data[i] = bytes

				} else {
					blockNumber := strconv.Itoa((nodeId % 3000) + int(block.Height*10))
					err = bc.db.Update(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte(blocksBucket))
						data[i] = b.Get([]byte(blockNumber))

						return nil
					})
				}
				log.Println("5")
			}
			ok, err := enc.Verify(data)
			if !ok {
				err = enc.Reconstruct(data)
			}
			for i := 0; i < 7; i++ {
				log.Println("^^^^^^")
				var result []byte
				for j := 0; j < len(data[i])-1; j++ {
					if data[i][j] != 0x20 {
						result = append(result, data[i][j])
					}
				}

				log.Println(DeserializeBlock(result))
			}

		}
	}
}
