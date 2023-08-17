package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net"
	"regexp"
	"strconv"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/klauspost/reedsolomon"
	"github.com/sleeg00/blockchain/proto"
	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

var mutex sync.Mutex

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var RS string
var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:3000", "localhost:3001", "localhost:3002", "localhost:3003", "localhost:3004", "localhost:3005",
	"localhost:3006", "localhost:3007", "localhost:3008", "localhost:3009"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction)
var keys []string
var lastIndex int

type server struct {
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
	defer bc.db.Close()
	UTXOSet := UTXOSet{Blockchain: bc}
	UTXOSet.Reindex()
	defer bc.db.Close()
	return &blockchain.CreateBlockchainResponse{
		Response: "Success",
	}, nil
}

func (s *server) Send(ctx context.Context, req *proto.SendRequest) (*proto.SendResponse, error) {

	log.Println("Send - Server receive a block")
	Tx := convertToTransaction(req.Block)

	bc := NewBlockchain(req.NodeTo)

	defer bc.db.Close()
	b := req.Block.PrevBlockHash
	UTXOSet := UTXOSet{Blockchain: bc}

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

	UTXOSet.Update(&block)

	response := &proto.SendResponse{
		Byte: b,
	}

	return response, nil
}

func (s *server) Mining(ctx context.Context, req *proto.MiningRequest) (*proto.MiningResponse, error) {

	var tx []Transaction
	for key := range mempool {
		tx = append(tx, mempool[key])
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

	log.Println("SendTrasaction - Server")
	tx := convertToOneTransaction(req.Transaction)

	mempool[hex.EncodeToString(req.Transaction.Id)] = tx

	for key := range mempool {
		keys = append(keys, key)
	}
	return &proto.ResponseTransaction{}, nil
}
func (s *server) CheckZombie(ctx context.Context, req *proto.CheckZombieRequest) (*proto.CheckZombieResponse, error) {
	return &proto.CheckZombieResponse{}, nil
}

// 0, 1, 2,3, 4, 5, 6, -- 7!
// 7번째 블록이 생성될 떄 RSEncoding을 진행하면 문제없이 이전 블럭 해쉬값을 알 수 있다.!!
func (s *server) RSEncoding(ctx context.Context, req *proto.RSEncodingRequest) (*proto.RSEncodingResponse, error) {
	nodeId, err := strconv.Atoi(req.NodeId)
	checkErr(err)

	bc := NewBlockchain(req.NodeId)
	defer bc.db.Close()
	bci := bc.Iterator()
	log.Println("NF: ", req.NF, "  F: ", req.F)
	enc, err := reedsolomon.New(7, 3) //비잔틴 장애 내성 가지도록 설계
	checkErr(err)
	//샤딩할 부분을 나눈다

	data := make([][]byte, int(10))

	for i := 6; i >= 0; i-- {
		block, err := bci.Next()
		checkErr(err)
		newBlockBytes := block.Serialize()
		// 비어있는 곳을 0으로 채운 후, newBlockBytes의 내용을 복사합니다
		data[i] = make([]byte, 1280)
		copy(data[i], newBlockBytes)

		for j := len(newBlockBytes); j < 1280; j++ {
			data[i][j] = 0x00
		}

	}

	for i := 7; i < 10; i++ {
		data[i] = make([]byte, 1280)
	}

	err = enc.Encode(data)

	checkErr(err)
	ok, err := enc.Verify(data)
	checkErr(err)
	log.Println("OOOKKK", ok)
	if ok == false {
		log.Panicln("!@#!@#!@#!@#!@#!@#!@")
	}
	bci = bc.Iterator()

	bci.Next()

	for i := 0; i < 7; i++ {
		block, err := bci.Next()
		checkErr(err)
		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			blockInDb := b.Get(block.Hash)

			if blockInDb == nil {
				return nil
			}

			err := b.Delete(block.Hash)
			if err != nil {
				log.Panic(err)
			}

			return nil
		})
		checkErr(err)
	}

	save := data[nodeId%3000]
	log.Println(data)
	end := 0
	check := false
	if req.Count != 0 {
		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			c := b.Cursor()
			if c != nil {
				keyBytes, _ := c.Last()

				targetBytes := []byte{72, 97, 115, 104}
				for keyBytes != nil {
					log.Println("keyBytes ", keyBytes)
					if len(keyBytes) >= 4 && bytes.HasPrefix(keyBytes[:4], targetBytes) {
						for j := 3; j < len(keyBytes); j++ {
							if keyBytes[j] == '~' {
								check = true
							} else if check {
								m := math.Pow(10, float64(len(keyBytes)-j-1))
								end += int(keyBytes[j]-48) * int(m)
							}
						}
						log.Println("keyBytes1 ", keyBytes)
						return nil
					}
					keyBytes, _ = c.Prev() // 이전 항목으로 이동
				}
			}
			return nil
		})
		checkErr(err)

		startIndex := "Hash" + strconv.Itoa(end+2) + "~" + strconv.Itoa(end+1+int(req.NF)+int(req.F)-1)
		log.Println(startIndex)
		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			b.Put([]byte(startIndex), save)

			return nil
		})
		checkErr(err)
	} else {

		startIndex := "Hash0~" + strconv.Itoa(int(req.NF)+(int(req.F)-1))
		log.Println("startIndex:", []byte(startIndex))
		log.Println("save------", save)
		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			b.Put([]byte(startIndex), save)

			return nil
		})
		checkErr(err)
	}
	log.Println(data)
	return &blockchain.RSEncodingResponse{}, nil
}

func (s *server) FindChunkTransaction(ctx context.Context, req *proto.FindChunkTransactionRequest) (*proto.FindChunkTransactionReponse, error) {

	bc := NewBlockchainRead(req.NodeId)
	defer bc.db.Close()

	var data [][]byte
	cnt := 0
	var foundTx Transaction // 구조체 변수 선언
	check := false
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		c := b.Cursor()
		if c != nil {
			// 키 공간을 역순으로 순회하여 가장 마지막 항목을 찾음
			keyBytes, _ := c.Last()

			targetBytes := []byte{72, 97, 115, 104}

			for keyBytes != nil {

				if len(keyBytes) >= 4 && bytes.HasPrefix(keyBytes[:4], targetBytes) {

					// value의 내용을 새로운 슬라이스에 복사하여 추가 (깊은 복사)
					value := b.Get(keyBytes)
					log.Println("value", len(value))

					data = append(data, make([]byte, len(value)))
					data[cnt] = make([]byte, len(value))
					copy(data[cnt], value)

					cnt++

					block := DeserializeBlock(data[cnt-1])

					for _, tx := range block.Transactions {
						if bytes.Equal(tx.ID, req.VinId) {
							log.Println("TX----", tx)

							for _, tx := range block.Transactions {
								if bytes.Equal(tx.ID, req.VinId) {
									check = true

									log.Println(req.NodeId, "에서 발견")
									foundTx = tx.TrimmedCopy() // 찾은 트랜잭션의 값을 복사하여 할당
									return nil
								}
							}

						}
					}
				}
				keyBytes, _ = c.Prev() // 이전 항목으로 이동

			}

		}

		return nil
	})
	checkErr(err)

	if check {
		return &proto.FindChunkTransactionReponse{
			Transaction: convertToProtoTransaction(foundTx),
		}, nil
	}
	return &proto.FindChunkTransactionReponse{
		Transaction: nil,
	}, nil
}
func (s *server) GetShard(ctx context.Context, req *proto.GetShardRequest) (*proto.GetShardResponse, error) {
	bc := NewBlockchainRead(req.NodeId)
	//defer bc.db.Close()

	var data [][]byte
	var lastKey string
	lastKey = ""
	log.Println("lastKEy!!!!", []byte(lastKey))

	cnt := 0
	data = make([][]byte, 1)
	check := false
	start := 0
	end := 0
	im := 2

	var list []int32
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		c := b.Cursor()
		if c != nil {
			// 키 공간을 역순으로 순회하여 가장 마지막 항목을 찾음
			keyBytes, _ := c.Last()

			lastKey = string(keyBytes)

			targetBytes := []byte{72, 97, 115, 104}

			for keyBytes != nil {

				if len(keyBytes) >= 4 && bytes.HasPrefix(keyBytes[:4], targetBytes) {
					log.Println("keyBytes:", keyBytes)
					// value의 내용을 새로운 슬라이스에 복사하여 추가 (깊은 복사)
					value := b.Get(keyBytes)

					data = append(data, make([]byte, len(value)))
					data[cnt] = value

					for i := len(keyBytes) - 1; i >= 4; i-- {
						var keyCheck int
						keyCheck = 0
						log.Println(keyCheck)
						if keyBytes[i] == '~' {
							check = true
							keyCheck = i
						} else if !check {
							log.Println("현재 KeyBytes", keyBytes)
							m := math.Pow(10, float64(len(keyBytes)-i-1))
							end += int(keyBytes[i]-48) * int(m)
						} else if check {
							log.Println("현재2 KeyBytes", keyBytes)
							m := math.Pow(10, float64(keyCheck-1-i))
							start += int(keyBytes[i]-48) * int(m)
						}
					}

					list = append(list, make([]int32, 2)...)

					list[im-2] = int32(start)
					list[im-1] = int32(end)
					cnt++
					im += 2
					check = false
					start = 0
					end = 0
					log.Println(list)

				}
				keyBytes, _ = c.Prev() // 이전 항목으로 이동
				lastKey = string(keyBytes)
			}

		}

		return nil
	})

	checkErr(err)

	return &proto.GetShardResponse{
		Bytes: data, // 깊은 복사된 슬라이스를 반환
		List:  list,
	}, nil
}

func (s *server) DeleteMempool(ctx context.Context, req *proto.DeleteMempoolRequest) (*proto.DeleteMempoolResponse, error) {
	log.Println("\n\n\n\nMEMPOLL SIZE", len(mempool))
	for key := range mempool {
		delete(mempool, key)
	}

	return &proto.DeleteMempoolResponse{}, nil
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
func convertToProtoTransaction(tx Transaction) *proto.Transaction {
	pbVin := make([]*proto.TXInput, len(tx.Vin))
	for i, vin := range tx.Vin {
		pbVin[i] = &proto.TXInput{
			Txid:      vin.Txid,
			Vout:      int64(vin.Vout),
			Signature: vin.Signature,
			PubKey:    vin.PubKey,
		}
	}

	pbVout := make([]*proto.TXOutput, len(tx.Vout))
	for i, vout := range tx.Vout {
		pbVout[i] = &proto.TXOutput{
			Value:      int64(vout.Value),
			PubKeyHash: vout.PubKeyHash,
		}
	}

	return &proto.Transaction{
		Id:   tx.ID,
		Vin:  pbVin,
		Vout: pbVout,
	}
}
func extractNumbersFromKey(key string) (int, int) {
	// 정규식 패턴을 사용하여 숫자를 추출
	re := regexp.MustCompile(`hash(\d+)~(\d+)`)
	matches := re.FindStringSubmatch(key)

	if len(matches) < 3 {
		// 매치되는 숫자가 없을 경우 오류 처리
		log.Fatal("Could not extract numbers from key")
	}

	start, _ := strconv.Atoi(matches[1])
	end, _ := strconv.Atoi(matches[2])

	return start, end
}
