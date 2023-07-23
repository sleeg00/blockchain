package main

import (
	"context"
	"log"

	"github.com/boltdb/bolt"
	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

func (cli *CLI) send(from, to string, amount int, node_id string, mineNow bool) {
	log.Println("Send")

	bc := NewBlockchainRead(node_id)
	UTXOSet := UTXOSet{bc}

	wallets, err := NewWallets(node_id) //wallet.node_id 확인
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := NewUTXOTransaction(&wallet, to, amount, &UTXOSet) //돈이 있는지 검사

	cbTx := NewCoinbaseTX(from, "") //마이닝했기 때문에 새로운 TX가 발생한다
	txs := []*Transaction{cbTx, tx}

	var lastHash []byte
	var lastHeight int

	for _, tx := range txs {
		// TODO: ignore transaction if it's not valid
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}

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

	blockDone := make(chan struct{})

	var newBlock *Block
	// 블록 생성 작업을 고루틴으로 실행
	go func() {
		newBlock = NewBlock(txs, lastHash, lastHeight+1)
		blockDone <- struct{}{} // 블록 생성 완료를 알리는 신호를 채널에 전송

	}()

	// 블록 생성이 완료될 때까지 대기
	<-blockDone

	// DB 닫기

	// 새로운 슬라이스를 만들고 txs의 값을 복사
	var protoTransactions []*blockchain.Transaction
	defer bc.db.Close()
	for _, tx := range txs {

		// 각 *Transaction을 proto.Transaction으로 매핑해서 protoTransactions 슬라이스에 추가합니다.
		protoTx := &blockchain.Transaction{
			Id:   tx.ID,
			Vin:  []*blockchain.TXInput{},
			Vout: []*blockchain.TXOutput{},
		}
		for _, vin := range tx.Vin {
			pbVin := &blockchain.TXInput{
				Txid:      vin.Txid,
				Vout:      int64(vin.Vout),
				Signature: vin.Signature,
				PubKey:    vin.PubKey,
			}
			protoTx.Vin = append(protoTx.Vin, pbVin)
		}
		for _, vout := range tx.Vout {
			pbVout := &blockchain.TXOutput{
				Value:      int64(vout.Value),
				PubKeyHash: vout.PubKeyHash,
			}
			protoTx.Vout = append(protoTx.Vout, pbVout)
		}
		protoTransactions = append(protoTransactions, protoTx)

	}

	block := &blockchain.Block{
		Timestamp:     newBlock.Timestamp,
		Transactions:  protoTransactions,
		PrevBlockHash: newBlock.PrevBlockHash,
		Hash:          newBlock.Hash,
		Nonce:         int32(newBlock.Nonce),
		Height:        int32(newBlock.Height),
	}

	for i := 0; i < len(knownNodes); i++ {

		if knownNodes[i][10:] == node_id {

			// 서버에 보낼 요청 메시지 생성
			request := &blockchain.SendRequest{
				From:     from,
				To:       to,
				Amount:   int32(amount),
				NodeFrom: node_id,
				NodeTo:   knownNodes[i][10:],
				MineNow:  mineNow,
				Block:    block,
			}
			log.Println("Request에러는 아님 ")
			// 서버에 요청 보내기
			response1, err1 := cli.blockchain.Send(context.Background(), request)
			if err1 != nil {
				log.Printf("Error sending request to node %s: %v", node_id, err)
				continue
			}

			// 서버 응답 처리...
			log.Printf("Received response from node %s: %v", node_id, response1)
		} else {
			conn, err := grpc.Dial(knownNodes[i], grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Failed to connect to gRPC server: %v", err)
			}
			defer conn.Close()

			client := blockchain.NewBlockchainServiceClient(conn)
			if err != nil {
				log.Printf("Error creating client for node %s: %v", knownNodes[i], err)
				continue
			}

			// 서버에 보낼 요청 메시지 생성
			request := &blockchain.SendRequest{
				From:     from,
				To:       to,
				Amount:   int32(amount),
				NodeFrom: node_id,
				NodeTo:   knownNodes[i][10:],
				MineNow:  mineNow,
				Block:    block,
			}

			// 서버에 요청 보내기
			response2, err2 := client.Send(context.Background(), request)
			if err2 != nil {
				log.Printf("Error sending request to node %s: %v", node_id, err)
				continue
			}

			// 서버 응답 처리...
			log.Printf("Received response from node %s: %v", node_id, response2)
		}

	}

}

/*
func (cli *CLI) send(from, to string, amount int, nodeID string, mineNow bool) {

	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(nodeID)

	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := NewUTXOTransaction(&wallet, to, amount, &UTXOSet)

	if mineNow {
		cbTx := NewCoinbaseTX(from, "") //마이닝했기 때문에 새로운 TX가 발생한다
		txs := []*Transaction{cbTx, tx} //두개의 TX들이 들어있는 []이다

		newBlock := bc.MineBlock(txs) //두개의 TX가 들어있는 블럭을 마이닝한다
		UTXOSet.Update(newBlock)
	} else {
		sendTx(knownNodes[0], tx)
	}

	fmt.Println("Success!")
}
*/
