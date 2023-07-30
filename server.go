package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/sleeg00/blockchain/proto"
	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

var mutex sync.Mutex

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:3000", "localhost:3001", "localhost:3002"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction)
var keys []string

type server struct {
}

type addr struct {
	AddrList []string
}

type block struct {
	AddrFrom string
	Block    []byte
}

type getblocks struct {
	AddrFrom string
}

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type tx struct {
	AddFrom     string
	Transaction []byte
}

type verzion struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}

func sendAddr(address string) {
	nodes := addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

func sendBlock(addr string, b *Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendTx(addr string, tnx *Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

// Version을 보낼 노드, Version을 보낸 노드의 Blockchain
func sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.GetBestHeight() //Version을 보낸 노드의 최고 길이는?
	payload := gobEncode(verzion{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

/*
bc := NewBlockchain(nodeID) //여기서 DB안 닫았음.

	for i := 0; i < len(knownNodes); i++ {
		if nodeAddress != knownNodes[i] {
			sendVersion(knownNodes[i], bc)
		}
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
*/
func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))
	requestBlocks()
}

func handleBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := UTXOSet{Blockchain: bc}
		UTXOSet.Reindex()
	}
}

func handleInv(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetBlocks(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleGetData(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(payload.AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

func handleTx(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*Transaction

			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := NewCoinbaseTX(miningAddress, "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := UTXOSet{Blockchain: bc}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range knownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}

func handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload verzion

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	log.Println(myBestHeight, " ", foreignerBestHeight)
	log.Println("모든 노드 + ", knownNodes)

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	sendAddr(payload.AddrFrom)
	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func handleConnection(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

func StartServer(nodeID, minerAddress string) {
	LocalNode := "localhost:" + nodeID
	srv := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(LocalNode))

	if err != nil {
		log.Fatalf("서버 연결 안됨")
	}
	blockchainService := &server{}
	blockchain.RegisterBlockchainServiceServer(srv, blockchainService)

	log.Println("Server listening on localhost:", nodeID)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
func (s *server) CreateWallet(ctx context.Context, req *blockchain.CreateWalletRequest) (*blockchain.CreateWalletResponse, error) {

	wallets, _ := NewWallets(req.NodeId)

	address := wallets.CreateWallet()

	wallets.SaveToFile(req.NodeId)

	fmt.Printf("Your new address: %s\n", address)
	return &blockchain.CreateWalletResponse{
		Address: address,
	}, nil
}

func (s *server) CreateBlockchain(ctx context.Context, req *blockchain.CreateBlockchainRequest) (*blockchain.CreateBlockchainResponse, error) {
	if !ValidateAddress(req.Address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := CreateBlockchain(req.Address, req.NodeId)

	UTXOSet := UTXOSet{Blockchain: bc}
	UTXOSet.Reindex()
	defer bc.db.Close()
	return &blockchain.CreateBlockchainResponse{
		Response: "Success",
	}, nil
}

func (s *server) Send(ctx context.Context, req *proto.SendRequest) (*proto.SendResponse, error) {

	log.Println("1")
	Tx := convertToTransaction(req.Block)

	bc := NewBlockchain(req.NodeTo)
	log.Println("2")
	defer bc.db.Close()
	b := req.Block.PrevBlockHash
	UTXOSet := UTXOSet{Blockchain: bc}
	log.Println("3")
	block := Block{
		Timestamp:     req.Block.Timestamp,
		PrevBlockHash: req.Block.PrevBlockHash,
		Transactions:  Tx,
		Hash:          req.Block.Hash,
		Nonce:         int(req.Block.Nonce),
		Height:        int(req.Block.Height),
	}
	log.Println("?")
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		a1 := block.Hash
		b1 := block.Serialize()
		err := b.Put(a1, b1)
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = block.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	log.Println("??")
	UTXOSet.Update(&block)

	response := &proto.SendResponse{
		Byte: b,
	}

	return response, nil
}

func (s *server) Mining(ctx context.Context, req *proto.MiningRequest) (*proto.MiningResponse, error) {

	var tx []Transaction
	for _, key := range keys {
		transaction := mempool[key]
		tx = append(tx, transaction)

	}
	log.Println("mempool에 저장한 TX들", tx)
	changeTx := convertToProtoTransactions(tx)

	response := &proto.MiningResponse{

		Response:     "Mining response 2",
		Transactions: changeTx,
	}
	return response, nil
}

func (s *server) FindMempool(ctx context.Context, req *proto.FindMempoolRequest) (*proto.FindMempoolResponse, error) {

	tx := makeTransactionNotPointer(mempool[req.HexTxId])
	return &proto.FindMempoolResponse{
		Transaction: tx,
	}, nil
}
func (s *server) SendTransaction(ctx context.Context, req *proto.SendTransactionRequest) (*proto.ResponseTransaction, error) {

	tx := convertToOneTransaction(req.Transaction)

	mempool[hex.EncodeToString(req.Transaction.Id)] = tx
	log.Println(len(mempool))
	log.Println("\n hex: ", hex.EncodeToString(req.Transaction.Id))
	for key := range mempool {
		keys = append(keys, key)
	}
	return &proto.ResponseTransaction{}, nil
}

func (s *server) SendBlock(ctx context.Context, req *proto.SendBlockRequest) (*proto.SendBlockResponse, error) {
	log.Println("블럭을 잘 전달받았음 ")
	bc := NewBlockchain(req.NodeId)
	UTXOSet := UTXOSet{Blockchain: bc}
	defer bc.db.Close()
	Tx := convertToTransaction(req.Block)

	block := Block{
		Timestamp:     req.Block.Timestamp,
		PrevBlockHash: req.Block.PrevBlockHash,
		Transactions:  Tx,
		Hash:          req.Block.Hash,
		Nonce:         int(req.Block.Nonce),
		Height:        int(req.Block.Height),
	}

	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(block.Hash, block.Serialize())
		log.Println()
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = block.Hash
		req.Block.PrevBlockHash = block.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXOSet.Reindex()
	return &proto.SendBlockResponse{
		Response: "Success",
	}, nil
}
func convertToTransaction(pbBlock *proto.Block) []*Transaction {
	var transactions []*Transaction
	for _, tx := range pbBlock.Transactions {
		transaction := &Transaction{
			ID:   tx.Id,
			Vin:  []TXInput{},
			Vout: []TXOutput{},
		}
		for _, pbVin := range tx.Vin {
			vin := TXInput{
				Txid:      pbVin.Txid,
				Vout:      int(pbVin.Vout),
				Signature: pbVin.Signature,
				PubKey:    pbVin.PubKey,
			}
			transaction.Vin = append(transaction.Vin, vin)
		}
		for _, pbVout := range tx.Vout {
			vout := TXOutput{
				Value:      int(pbVout.Value),
				PubKeyHash: pbVout.PubKeyHash,
			}
			transaction.Vout = append(transaction.Vout, vout)
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

func convertToOneTransaction(tx *proto.Transaction) Transaction {
	var transactions Transaction

	transaction := &Transaction{
		ID:   tx.Id,
		Vin:  []TXInput{},
		Vout: []TXOutput{},
	}
	for _, pbVin := range tx.Vin {
		vin := TXInput{
			Txid:      pbVin.Txid,
			Vout:      int(pbVin.Vout),
			Signature: pbVin.Signature,
			PubKey:    pbVin.PubKey,
		}
		transaction.Vin = append(transaction.Vin, vin)
	}
	for _, pbVout := range tx.Vout {
		vout := TXOutput{
			Value:      int(pbVout.Value),
			PubKeyHash: pbVout.PubKeyHash,
		}
		transaction.Vout = append(transaction.Vout, vout)
	}
	transactions = *transaction

	return transactions
}
func convertToProtoTransactions(txList []Transaction) []*proto.Transaction {
	var protoTxs []*proto.Transaction

	for _, tx := range txList {
		protoTxs = append(protoTxs, &proto.Transaction{
			Id:   tx.ID,
			Vin:  convertToProtoInputs(tx.Vin),
			Vout: convertToProtoOutputs(tx.Vout),
		})
	}

	return protoTxs
}

func convertToProtoInputs(inputs []TXInput) []*proto.TXInput {
	var protoInputs []*proto.TXInput
	for _, input := range inputs {
		protoInputs = append(protoInputs, &proto.TXInput{
			Txid:      input.Txid,
			Vout:      int64(input.Vout),
			Signature: input.Signature,
			PubKey:    input.PubKey,
		})
	}
	return protoInputs
}

func convertToProtoOutputs(outputs []TXOutput) []*proto.TXOutput {
	var protoOutputs []*proto.TXOutput
	for _, output := range outputs {
		protoOutputs = append(protoOutputs, &proto.TXOutput{
			Value:      int64(output.Value),
			PubKeyHash: output.PubKeyHash,
		})
	}
	return protoOutputs
}
