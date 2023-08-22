package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

func RsEncoding(count int32, f int32, NF int32) {
	var wg sync.WaitGroup
	for i := 0; i < len(knownNodes); i++ {
		go func(i int) {
			serverAddress := fmt.Sprintf("localhost:%s", knownNodes[i][10:])
			conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
			if err != nil {
				log.Printf("Failed to connect to gRPC server: %v", err)
				return
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
		}(i)
	}
	wg.Wait()
	log.Println("1")
	/*
		failNodeslist := make([]string, 10)
		failNodesCheckCount := 0
			monitorInterval := time.Second * 10 // 10초마다 노드 상태를 모니터링

			resultCh := make(chan string)

			for {
				var wg sync.WaitGroup

				for _, nodeID := range knownNodes {
					wg.Add(1)
					go func(nodeID string) {
						defer wg.Done()
						failNodes, failNodesCheck := monitorNode(nodeID)
						if failNodesCheck == 1 {
							failNodeslist = append(failNodeslist, failNodes)
							failNodesCheckCount++
						}

					}(nodeID)
				}
				var data [][]byte
				var list []int32
				cnt := 0
				data = make([][]byte, 10)
				enc, err := reedsolomon.New(7, 3)
				checkErr(err)
				if failNodesCheckCount > 0 { //청크들 모음
					for k := 0; k < len(knownNodes); k++ {

						if knownNodes[k][10:] != failNodeslist[0] {
							serverAddress := fmt.Sprintf("localhost:%s", knownNodes[k][10:])

							conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

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
							if err != nil {
								log.Println("연결실패!", knownNodes[k])
								failNodes = append(failNodes, knownNodes[k][10:])
								failNodesCheck++
								continue
							}
							bytes := response.Bytes
							list = response.List
							log.Println("list: ", list)
							//log.Println(DeserializeBlock(bytes))

							size := len(bytes)

							cnt = 0
							for j := 0; ; j++ {
								if cnt == size-1 {
									break
								}
								data[cnt*10+k] = bytes[cnt]
								cnt++
							}
						}
					}
				}
				ok, err := enc.Verify(data)
				checkErr(err)
				log.Println(ok)

				enc.ReconstructData(data)
				wg.Wait() // 현재 라운드의 모든 고루틴이 종료될 때까지 대기합니다.

				time.Sleep(monitorInterval) // 10초 대기
			}

			close(resultCh)

			// 결과 채널에서 노드 상태 출력하기
			for result := range resultCh {
				fmt.Println(result)
			}
	*/
}

/*
func monitorNode(nodeID string) (failNodes string, failNodesCheck int32) {

	serverAddress := fmt.Sprintf("localhost:%s", nodeID)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

	log.Println(serverAddress)

	defer conn.Close()
	client := blockchain.NewBlockchainServiceClient(conn)
	cli := CLI{
		nodeID:     nodeID,
		blockchain: client,
	}
	request := &blockchain.CheckZombieRequest{}

	response, err := cli.blockchain.CheckZombie(context.Background(), request)
	log.Println(response)
	if err != nil {
		return nodeID, 1
	}
	return "Success", 0
}
*/
