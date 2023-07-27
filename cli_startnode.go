package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/sleeg00/blockchain/proto"
	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

func (cli *CLI) startNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if ValidateAddress(minerAddress) {
			log.Println("모든 노드 mempool에게서 TX를 가져와서 채굴노드 mempool에저장한다")
			var wg sync.WaitGroup //고루틴이 완료되기를 기다리기 위한 준비를 합니다.
			for i := 0; i < len(knownNodes); i++ {

				if knownNodes[i][10:] != nodeID {
					serverAddress := fmt.Sprintf("localhost:%s", knownNodes[i][10:])

					conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
					if err != nil {
						log.Fatalf("Failed to connect to gRPC server: %v", err)
					}
					defer conn.Close()

					client := blockchain.NewBlockchainServiceClient(conn)
					cli := CLI{
						nodeID:     nodeID,
						blockchain: client,
					}

					request := &proto.MiningRequest{
						NodeTo: knownNodes[i][10:],
					}

					response, err := cli.blockchain.Mining(context.Background(), request)
					if err != nil {
						log.Fatalf("Error during Mining API call: %v", err)
					}
					for _, protoTx := range response.Transactions {
						convertTx := convertFromProtoTransaction(protoTx)
						mempool[hex.EncodeToString(convertTx.ID)] = *convertTx

					}
				}

			}

			log.Println("모든 노드 mempool에게서 TX를 가져와서 채굴노드 mempool에저장을 끝냈다")

			//-------블럭 생성
			log.Println("블럭을 생성한다.")
			bc := NewBlockchainRead(nodeID)
			defer bc.db.Close()
			wallets, err := NewWallets(nodeID) //wallet.node_id 확인
			if err != nil {
				log.Panic(err)
			}
			wallet := wallets.GetWallet(minerAddress)
			log.Println(wallet, "이란 지갑이 존재한다~")

			var txs []*Transaction

			for id := range mempool {

				tx := mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
				log.Println(tx)
			}
			log.Println("내가 모은 mempool에 TX찍어보기", txs)
			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := NewCoinbaseTX(minerAddress, "") //->여기가 문제
			txs = append(txs, cbTx)

			var lastHash []byte
			var lastHeight int

			err = bc.db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(blocksBucket))
				lastHash = b.Get([]byte("l"))

				blockData := b.Get(lastHash)
				block := DeserializeBlock(blockData)

				lastHeight = block.Height

				return nil
			})
			if err != nil {
				log.Panic(err)
			}

			wg.Add(1)

			var newBlock *Block
			// 블록 생성 작업을 고루틴으로 실행
			go func() {
				newBlock = NewBlock(txs, lastHash, lastHeight+1)
				wg.Done() // 작업이 끝날 때 Done()을 호출하여 종료를 알림
			}()

			wg.Wait() // 블록 생성 완료를 기다림

			log.Println("모든 노드의 Mempool을 기반으로 채굴자노드에서 블럭을 생성을 끝냈다.")

			log.Println("모든 노드에게 블럭을 전달한다.")
			protoTx := makeClientTransaction(newBlock) //요청을 보낼 Block으로 마샬링한다

			for i := 0; i < len(knownNodes); i++ {
				wg.Add(1) // 고루틴의 수를 증가시킵니다.
				go func(index int) {
					defer wg.Done() // 해당 고루틴이 끝나면 WaitGroup에서 하나를 차감합니다.

					if nodeID == knownNodes[index][10:] {
						block := &proto.Block{
							Timestamp:     newBlock.Timestamp,
							Transactions:  protoTx,
							PrevBlockHash: newBlock.PrevBlockHash,
							Hash:          newBlock.Hash,
							Nonce:         int32(newBlock.Nonce),
							Height:        int32(newBlock.Height),
						}
						request := &proto.SendBlockRequest{
							NodeId: knownNodes[index][10:],
							Block:  block,
						}

						response, err := cli.blockchain.SendBlock(context.Background(), request)
						if err != nil {
							log.Fatalf("Error during Mining API call: %v", err)
						}

						if response.Response == "Success" {
							fmt.Println("Success")
						} else {
							fmt.Println("fail")
						}
					} else {
						serverAddress := fmt.Sprintf("localhost:%s", knownNodes[index][10:])

						conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
						if err != nil {
							log.Fatalf("Failed to connect to gRPC server: %v", err)
						}
						defer conn.Close()

						block := &proto.Block{
							Timestamp:     newBlock.Timestamp,
							Transactions:  protoTx,
							PrevBlockHash: newBlock.PrevBlockHash,
							Hash:          newBlock.Hash,
							Nonce:         int32(newBlock.Nonce),
							Height:        int32(newBlock.Height),
						}

						client := blockchain.NewBlockchainServiceClient(conn)
						cli := CLI{
							nodeID:     nodeID,
							blockchain: client,
						}

						request := &proto.SendBlockRequest{
							NodeId: knownNodes[index][10:],
							Block:  block,
						}

						response, err := cli.blockchain.SendBlock(context.Background(), request)
						if err != nil {
							log.Fatalf("Error during Mining API call: %v", err)
						}

						if response.Response == "Success" {
							fmt.Println("Success")
						} else {
							fmt.Println("fail")
						}
					}
				}(i)

			}
			wg.Wait()

			//for id := range mempool {
			//tx := mempool[id]
			//delete(mempool, tx.String())

			//}
			log.Println("모든 노드에게 블럭을 전달을 끝냈다.")
		} else {
			log.Panic("Wrong miner address!")
		}
	} else {
		StartServer(nodeID, minerAddress)
	}
}
func makeClientTransaction(newBlock *Block) []*proto.Transaction {
	var protoTransactions []*proto.Transaction

	for _, tx := range newBlock.Transactions {

		// 각 *Transaction을 proto.Transaction으로 매핑해서 protoTransactions 슬라이스에 추가합니다.
		protoTx := &proto.Transaction{
			Id:   tx.ID,
			Vin:  []*proto.TXInput{},
			Vout: []*proto.TXOutput{},
		}
		for _, vin := range tx.Vin {
			pbVin := &proto.TXInput{
				Txid:      vin.Txid,
				Vout:      int64(vin.Vout),
				Signature: vin.Signature,
				PubKey:    vin.PubKey,
			}
			protoTx.Vin = append(protoTx.Vin, pbVin)
		}
		for _, vout := range tx.Vout {
			pbVout := &proto.TXOutput{
				Value:      int64(vout.Value),
				PubKeyHash: vout.PubKeyHash,
			}
			protoTx.Vout = append(protoTx.Vout, pbVout)
		}
		protoTransactions = append(protoTransactions, protoTx)

	}
	return protoTransactions
}
