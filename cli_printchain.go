package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

func printChain(nodeID string) {

	bc := NewBlockchainRead(nodeID)

	defer bc.db.Close()

	bci := bc.Iterator()

	data := make([][]byte, 20)

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

	cnt := 0
	var failNodes = []string{}
	var failNodesCheck int
	var list []int32

	for k := 0; k < len(knownNodes); k++ {

		serverAddress := fmt.Sprintf("localhost:%s", knownNodes[k][10:])

		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
		checkErr(err)
		log.Println(serverAddress)
		defer conn.Close()
		client := blockchain.NewBlockchainServiceClient(conn)
		cli := CLI{
			nodeID:     knownNodes[k][10:],
			blockchain: client,
		}
		log.Println(cli)

		//여기서 지금 멈춤

		request := &blockchain.GetShardRequest{
			NodeId: knownNodes[k][10:],
			Height: int32(count),
		}

		response, err := cli.blockchain.GetShard(context.Background(), request)
		if err != nil {
			log.Println("연결실패!", knownNodes[k])

		} else {

			bytes := response.Bytes

			if k == 0 {
				list = response.List
			}

			//log.Println(DeserializeBlock(bytes))

			size := len(bytes)
			log.Println("SIZE", size)
			cnt = 0
			log.Println("List", list)
			for j := 0; ; j++ {
				if cnt == size {
					break
				}
				data[cnt*10+k] = bytes[cnt]

				cnt++

			}
		}

	}
	log.Println("여기")

	log.Println(list)
	for k := 0; k < len(knownNodes); k++ {

		serverAddress := fmt.Sprintf("localhost:%s", knownNodes[k][10:])

		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

		log.Println(serverAddress)
		defer conn.Close()
		if err != nil {

		} else {

			client := blockchain.NewBlockchainServiceClient(conn)
			cli := CLI{
				nodeID:     knownNodes[k][10:],
				blockchain: client,
			}
			request := &blockchain.CheckRsEncodingRequest{
				Bytes:  data,
				NodeId: knownNodes[k][10:],
			}

			response, err := cli.blockchain.CheckRsEncoding(context.Background(), request)
			if err != nil {
				log.Println("연결실패!", knownNodes[k])
				failNodes = append(failNodes, knownNodes[k][10:])
				failNodesCheck++
				data[k] = nil

			} else {
				check := response.Check
				//log.Println(DeserializeBlock(bytes))
				if check == true {
					log.Println(knownNodes[k], "는 값이 같다")
				} else {
					log.Println(knownNodes[k], "는 값 같지 않다")
				}

			}
		}
	}

	for k := 0; k < len(knownNodes); k++ {

		serverAddress := fmt.Sprintf("localhost:%s", knownNodes[k][10:])

		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

		log.Println(serverAddress)
		if err != nil {

		} else {
			defer conn.Close()
			client := blockchain.NewBlockchainServiceClient(conn)
			cli := CLI{
				nodeID:     knownNodes[k][10:],
				blockchain: client,
			}
			request := &blockchain.CheckRsEncodingRequest{
				Bytes:  data,
				NodeId: knownNodes[k][10:],
			}

			response, err := cli.blockchain.CheckRsEncoding(context.Background(), request)
			if err != nil {
				log.Println("연결실패!", knownNodes[k])
				failNodes = append(failNodes, knownNodes[k][10:])
				failNodesCheck++
				continue
			} else {
				check := response.Check
				//log.Println(DeserializeBlock(bytes))
				if check == true {
					log.Println(knownNodes[k], "는 값이 같다")
				} else {
					log.Println(knownNodes[k], "는 값 같지 않다")
				}

			}
		}
	}

	check := false
	listCheck := 4
	for k := 0; k < 4; k += 2 {
		for x := 0; x <= 20; x++ {
			log.Println(list[0] + list[1])
			var result []byte

			if x%7 == 0 && x != 0 {
				check = true

			} else if list[listCheck-2]+list[listCheck-1]+1 == int32(x) {

				check = false
				listCheck -= 2
			}
			log.Println(x)
			if !check && x != 9 {
				if data[x] == nil { //복구 해야할 때
					log.Println("XX", x)

					for y := 0; y < len(data[x]); y++ {
						result = append(result, data[x][y])
					}

					block := DeserializeBlock(result)
					log.Println("7")
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

				} else {
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
	}
}
