package main

import (
	"context"
	"log"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/sleeg00/blockchain/proto"
)

type SafeBlock struct {
	block *proto.Block
	mutex sync.Mutex
}

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
	log.Println("이전 해쉬-2 : ", block.PrevBlockHash)

	// 서버에 요청 보내기
	response1, err1 := cli.blockchain.Send(context.Background(), request)
	if err1 != nil {
		log.Println("Error sending request to node %s: %v", node_id, err1)

	}

	log.Println("이전 해쉬-3 : ", block.PrevBlockHash)
	// 서버 응답 처리...
	log.Printf("Received response from node %s: %v", node_id, response1)
	log.Println("이전 해쉬-4 : ", block.PrevBlockHash)

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

func send(from, to string, amount int, node_id string, mineNow bool) *Block {
	bc := NewBlockchainRead(node_id)

	UTXOSet := UTXOSet{Blockchain: bc}

	wallets, err := NewWallets(node_id) // wallet.node_id 확인
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := NewUTXOTransaction(&wallet, to, amount, &UTXOSet) // 돈이 있는지 검사

	cbTx := NewCoinbaseTX(from, "") // 마이닝했기 때문에 새로운 TX가 발생한다
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

	newBlock := NewBlock(txs, lastHash, lastHeight+1)

	return newBlock
}

func sendTrsaction(from, to string, amount int, node_id string, mineNow bool) *Transaction {
	bc := NewBlockchain(node_id)
	defer bc.db.Close()
	UTXOSet := UTXOSet{Blockchain: bc}

	wallets, err := NewWallets(node_id) //wallet.node_id 확인
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := NewUTXOTransaction(&wallet, to, amount, &UTXOSet) //돈이 있는지 검사

	// TODO: ignore transaction if it's not valid
	if !bc.VerifyTransaction(tx) {
		log.Panic("ERROR: Invalid transaction")
	}
	UTXOSet.UpdateTx(tx)
	return tx
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
