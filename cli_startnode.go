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
			log.Println("채굴자 노드 mempool에게서 TX를 가져와서 채굴노드 mempool에저장한다")

			request := &proto.MiningRequest{
				NodeTo: nodeID,
			}

			response, err := cli.blockchain.Mining(context.Background(), request)

			if err != nil {
				log.Fatalf("Error during Mining API call: %v", err)
			}
			for _, protoTx := range response.Transactions {
				convertTx := convertFromProtoTransaction(protoTx)

				mempool[hex.EncodeToString(convertTx.ID)] = convertTx
				log.Println(convertTx)
			}

			log.Println("모든 노드 mempool에게서 TX를 가져와서 채굴노드 mempool에저장을 끝냈다")
			log.Println("\n\n\n\n MEMpool size", len(mempool))
			if len(mempool) >= 2 && len(minerAddress) > 0 {

			MineTransactions:
				//-------블럭 생성
				log.Println("블럭을 생성한다.")

				bc := NewBlockchainRead(nodeID)

				wallets, err := NewWallets(nodeID) //wallet.node_id 확인
				if err != nil {
					log.Panic(err)
				}
				wallet := wallets.GetWallet(minerAddress)
				log.Println(wallet, "이란 지갑이 존재한다~")

				var txs []*Transaction

				for id := range mempool {

					tx := mempool[id]
					txs = append(txs, &tx)

				}

				if len(txs) == 0 {
					fmt.Println("All transactions are invalid! Waiting for new ones...")
					return
				}

				cbTx := NewCoinbaseTX(minerAddress, "") //채굴 보상자 트랜잭션 생성
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

				blockChannel := make(chan Block) // 채널 생성
				go func() {
					newblock := NewBlock(txs, lastHash, lastHeight+1)
					blockChannel <- newblock // 새 블록을 채널에 전달
				}()
				newblock := <-blockChannel // 채널로부터 결과를 받을 때까지 기다립니다.

				log.Println("모든 노드의 Mempool을 기반으로 채굴자노드에서 블럭을 생성을 끝냈다.")

				log.Println("모든 노드에게 블럭을 전달한다.")
				protoTransactions := makeClientTransactions(newblock.Transactions) //요청을 보낼 Block으로 마샬링한다
				log.Println(protoTransactions)
				/*
					for i := 0; i < len(knownNodes); i++ {

						if nodeID == knownNodes[i][10:] {
							bc.db.Close()
						}
						block := &proto.Block{
							Timestamp:     newblock.Timestamp,
							Transactions:  protoTransactions,
							PrevBlockHash: newblock.PrevBlockHash,
							Hash:          newblock.Hash,
							Nonce:         int32(newblock.Nonce),
							Height:        int32(newblock.Height),
						}

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

						byte := cli.gRPCsendBlockRequest(knownNodes[i][10:], block)
						newblock.PrevBlockHash = byte

					}
				*/
				var wg sync.WaitGroup
				for i := 0; i < len(knownNodes); i++ {
					wg.Add(1)
					go func(i int) {
						defer wg.Done()
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

						// 서버에 보낼 요청 메시지 생성
						request := &proto.DeleteMempoolRequest{}

						// 서버에 요청 보내기
						response, err := cli.blockchain.DeleteMempool(context.Background(), request)
						if err != nil {
							log.Panic(err)
						}
						log.Println(response)

					}(i)
				}
				wg.Wait()
				for key := range mempool {
					delete(mempool, key)
				}
				log.Println("\n\n\n\n삭제 후 Mempool 사이즈", len(mempool))
				if len(mempool) > 0 {
					goto MineTransactions
				}
				log.Println("모든 노드에게 블럭을 전달을 끝냈다.")
			}
		} else {
			log.Panic("Wrong miner address!")
		}
	} else {
		StartServer(nodeID, minerAddress)
	}
}
func makeClientTransactions(Transactions []*Transaction) []*proto.Transaction {
	var protoTransactions []*proto.Transaction

	for _, tx := range Transactions {

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

func makeOneTransaction(tx *Transaction) *proto.Transaction {

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

	return protoTx
}

func makeTransactionNotPointer(tx Transaction) *proto.Transaction {

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

	return protoTx
}
