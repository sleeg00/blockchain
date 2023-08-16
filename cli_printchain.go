package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/klauspost/reedsolomon"
	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

func (cli *CLI) printChain(nodeID string) {

	bc := NewBlockchainRead(nodeID)

	defer bc.db.Close()

	bci := bc.Iterator()

	data := make([][]byte, 10)

	var Height int
	for k := 0; k < len(knownNodes); k++ {

		serverAddress := fmt.Sprintf("localhost:%s", knownNodes[k][10:])

		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

		log.Println(serverAddress)

		defer conn.Close()
		client := blockchain.NewBlockchainServiceClient(conn)
		cli := CLI{
			nodeID:     knownNodes[k][10:],
			blockchain: client,
		}
		request := &blockchain.CheckZombieRequest{}

		response, err := cli.blockchain.CheckZombie(context.Background(), request)
		log.Println(response)
		if err != nil {
			log.Println(knownNodes[k], "에 연결 실패!")
			failNodes = append(failNodes, knownNodes[k][10:])
			failNodesCheck++
		}
	}

	log.Println("FailNode?!", failNodesCheck)
	invaildNodeCount := 10 - failNodesCheck
	f = (invaildNodeCount - 1) / 3
	NF = invaildNodeCount - f
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

		if len(block.PrevBlockHash) == 0 || x == Height%NF {
			log.Println("xxxxxxxx", x)
			break
		}
	}

	count := Height / NF

	cnt := 0
	var failNodes = []string{}
	var failNodesCheck int
	var list []int32
	for k := 0; k < len(knownNodes); k++ {

		serverAddress := fmt.Sprintf("localhost:%s", knownNodes[k][10:])

		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

		log.Println(serverAddress)
		if err != nil {
			log.Println(knownNodes[k], "에 연결 실패!")
			failNodes = append(failNodes, knownNodes[k][10:])
			failNodesCheck++
			continue
		} else {
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
				log.Println("cnt:", cnt*10+k)

				cnt++
			}
		}
	}

	check := false
	listCheck := len(list)
	log.Println("1")
	log.Println("NF", NF, "F", f)

	log.Println("2")
	restoreData := make([][]byte, len(data))

	for x := 0; x < int(list[0]+list[1]); x++ {

		var result []byte

		if int(list[listCheck-1])-f+1 == x && x != 0 {
			check = true

		} else if list[listCheck-2]+list[listCheck-1]+1 == int32(x) {
			check = false
			listCheck -= 2
		}

		if !check {
			if data[x] == nil { //복구 해야할 때
				log.Println("XX", x)
				enc, err := reedsolomon.New(7, 3)
				ok, err := enc.Verify(data)
				checkErr(err)
				log.Println(ok)
				err = enc.ReconstructData(data)
				checkErr(err)
				log.Println(restoreData[1])
				checkErr(err)

				for y := 0; y < len(restoreData[x]); y++ {

					result = append(result, restoreData[x][y])

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
