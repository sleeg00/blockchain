package main

import (
	"context"
	"log"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/sleeg00/blockchain/proto"
)

func (cli *CLI) request(from, to string, amount int, node_id string, mineNow bool, block *proto.Block, node_from string) {

	log.Println("이전블럭해시", block.PrevBlockHash)

	// 서버에 보낼 요청 메시지 생성
	request := &proto.SendRequest{
		From:     from,
		To:       to,
		Amount:   int32(amount),
		NodeFrom: node_from,
		NodeTo:   node_id,
		MineNow:  mineNow,
		Block:    block,
	}

	// 서버에 요청 보내기
	response1, err1 := cli.blockchain.Send(context.Background(), request)
	if err1 != nil {
		log.Panic("Error sending request to node %s: %v", node_id, err1)

	}

	// 서버 응답 처리...
	log.Printf("Received response from node %s: %v", node_id, response1)

}

func (cli *CLI) requestTransaction(from, to string, amount int, node_id string, mineNow bool, transaction *proto.Transaction, node_from string) {
	// 서버에 보낼 요청 메시지 생성
	request := &proto.SendTransactionRequest{
		From:        from,
		To:          to,
		Amount:      int32(amount),
		NodeFrom:    node_from,
		NodeTo:      node_id,
		MineNow:     mineNow,
		Transaction: transaction,
	}

	log.Println("Request에러는 아님 ")
	// 서버에 요청 보내기
	response1, err1 := cli.blockchain.SendTransaction(context.Background(), request)
	if err1 != nil {
		log.Panic("Error sending request to node %s: %v", node_id, err1)

	}

	// 서버 응답 처리...
	log.Printf("Received response from node %s: %v", node_id, response1)

}

func send(from, to string, amount int, node_id string, mineNow bool) Block {

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

	var wg sync.WaitGroup
	wg.Add(1)

	var newBlock Block
	// 블록 생성 작업을 고루틴으로 실행
	go func() {
		newBlock = NewBlockNotAddress(txs, lastHash, lastHeight+1)
		wg.Done() // 작업이 끝날 때 Done()을 호출하여 종료를 알림
	}()
	wg.Wait() // 블록 생성 완료를 기다림

	return newBlock
}

func sendTrsaction(from, to string, amount int, node_id string, mineNow bool) *proto.Transaction {
	bc := NewBlockchainRead(node_id)
	defer bc.db.Close()
	UTXOSet := UTXOSet{bc}

	wallets, err := NewWallets(node_id) //wallet.node_id 확인
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := NewUTXOTransaction(&wallet, to, amount, &UTXOSet) //돈이 있는지 검사

	cbTx := NewCoinbaseTX(from, "") //마이닝했기 때문에 새로운 TX가 발생한다
	txs := []*Transaction{cbTx, tx}

	for _, tx := range txs {
		// TODO: ignore transaction if it's not valid
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	var protoTransactions *proto.Transaction

	for _, tx := range txs {

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
		protoTransactions = protoTx

	}

	return protoTransactions
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
