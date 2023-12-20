package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/klauspost/reedsolomon"
	"github.com/sleeg00/blockchain/proto"
	blockchain "github.com/sleeg00/blockchain/proto"
	"github.com/xuri/excelize/v2"
	"google.golang.org/grpc"
)

var mutex sync.Mutex

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var list []int32
var RS string
var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:3000", "localhost:3001", "localhost:3002", "localhost:3003", "localhost:3004", "localhost:3005",
	"localhost:3006", "localhost:3007", "localhost:3008", "localhost:3009"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction)
var keys []string
var lastIndex int
var rs_server_process [100000][10]string
var cnt3 int

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
	cnt2 = 0
	cnt3 = 0
	done := make(chan struct{}) // done 채널 생성
	LocalNode := "localhost:" + nodeID
	srv := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(LocalNode))

	if err != nil {
		log.Fatalf("서버 연결 안됨")
	}
	blockchainService := &server{}
	blockchain.RegisterBlockchainServiceServer(srv, blockchainService)

	log.Println("Server listening on localhost:", nodeID)

	cpuUsageTicker := time.NewTicker(1 * time.Second) // CPU 사용률을 확인할 주기를 설정합니다. (예: 1초마다)
	defer cpuUsageTicker.Stop()

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
		close(done) // 서버가 종료되면 done 채널을 닫음
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM) // Ctrl + C 및 기타 종료 시그널을 처리
	for {
		select {
		case <-done:

		case sig := <-stop:
			log.Printf("Received signal: %v", sig)
			startTime := time.Now()

			// 서버 종료 후에 로직 실행
			log.Println("Server stopped. Executing reconnect and send request logic...")

			conn, err := grpc.Dial("localhost:3010", grpc.WithInsecure())
			if err != nil {
				log.Printf("Failed to connect to 3010 gRPC server: %v", err)
				return
			}
			defer conn.Close()

			newClient := blockchain.NewBlockchainServiceClient(conn)
			var hash []string
			for i := 99; i >= 0; i-- {
				hash = append(hash, "Hash"+strconv.Itoa(i*7)+"~"+strconv.Itoa((i*7)+6))
			}
			req := &blockchain.DataRequest{
				Hash:   hash,
				NodeId: nodeID,
			}

			res, err := newClient.DataTransfer(context.Background(), req)

			if err != nil {
				log.Printf("SendData failed: %v", err)
				return
			}

			if res.Success {
				log.Println("Data send successfully")
			} else {
				log.Println("Data send failed")
			}

			log.Println("Reconnect and send request logic executed")
			elapsedTime := time.Since(startTime)
			fmt.Printf("Total time taken: %s\n", elapsedTime)
			return
		case <-cpuUsageTicker.C:
			// 주기적으로 CPU 사용률을 확인하고 출력합니다.
			cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", os.Getpid()), "-o", "%cpu")
			output, err := cmd.Output()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			cpuUsageStr := strings.TrimSpace(string(output))
			fmt.Println("Current CPU Usage:", cpuUsageStr)
			cnt2++
			rs_server[cnt2][3] = cpuUsageStr
			rs_server_process[cnt3][0] = cpuUsageStr
			cnt3++
		}
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
func (s *server) SendServer(ctx context.Context, req *proto.SendRequest) (*proto.SendResponse, error) {
	log.Println("Send - Server receive a block")
	log.Println("1")

	log.Println("2")

	log.Println("3")

	log.Println("4")

	return &proto.SendResponse{}, nil
}

var block Block

func (s *server) Connect(ctx context.Context, req *proto.ConnectServerRequest) (*proto.ConnectServerResponse, error) {
	//	log.Println("Connect Block")
	Tx := convertToTransaction(req.Block)

	b := req.Block.PrevBlockHash

	block = Block{
		Timestamp:     req.Block.Timestamp,
		PrevBlockHash: req.Block.PrevBlockHash,
		Transactions:  Tx,
		Hash:          req.Block.Hash,
		Nonce:         int(req.Block.Nonce),
		Height:        int(req.Block.Height),
	}

	return &proto.ConnectServerResponse{Byte: b}, nil
}

func (s *server) Send(ctx context.Context, req *proto.SendRequest) (*proto.SendResponse, error) {
	//	log.Println("Send - Server receive a block")

	bc := NewBlockchain(req.NodeTo)
	defer bc.db.Close()

	UTXOSet := UTXOSet{Blockchain: bc}
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
	//	log.Println("블럭의 높이는:", block.Height)
	return &proto.SendResponse{}, nil
}
func (s *server) SaveFileSystem(ctx context.Context, req *proto.SendRequest) (*proto.SendResponse, error) {
	bc := NewBlockchain(req.NodeTo)
	defer bc.db.Close()
	for {
		count := 0

		bci := bc.Iterator()
		nodeId, err := strconv.Atoi(req.NodeTo)
		var checkList []int
		enc, err := reedsolomon.New(7, 3) //비잔틴 장애 내성 가지도록 설계
		checkErr(err)
		//샤딩할 부분을 나눈다

		data := make([][]byte, int(10))

		for i := 0; ; i++ {
			block, err := bci.Next()

			checkErr(err)
			if block.Height%7 == 0 {
				break
			}
		}
		log.Println("0")
		for i := 6; i >= 0; i-- {
			block, err := bci.Next() // 6
			log.Println(block.Height)
			checkErr(err)

			log.Println("RSEncoding Block Height", block.Height)
			newBlockBytes := block.Serialize()

			data[i] = make([]byte, 2048)
			copy(data[i], newBlockBytes)
			checkList = append(checkList, block.Height)
		}

		for i := 7; i < 10; i++ {
			data[i] = make([]byte, 2048)
		}
		err = enc.Encode(data)
		log.Println("1")
		checkErr(err)
		log.Println("2")
		ok, err := enc.Verify(data)
		checkErr(err)
		log.Println("3")

		if ok == false {
			log.Panicln("!@#!@#!@#!@#!@#!@#!@")
		}
		bci = bc.Iterator()

		for i := 0; ; i++ {
			block, err := bci.Next()
			checkErr(err)
			if block.Height%7 == 0 {
				break
			}
		}

		for i := 0; i < 7; i++ {
			block, err := bci.Next() //6
			log.Println("block:", block.Height)
			log.Println("checkList: ", checkList)
			checkErr(err)
			if block.Height == checkList[i] {
				err = bc.db.Update(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte(blocksBucket))
					blockInDb := b.Get(block.Hash)

					if blockInDb == nil {
						return nil
					}
					log.Println("삭제한 블럭 Height", DeserializeBlock(blockInDb).Height)
					err := b.Delete(block.Hash)
					err = b.Delete(bc.tip)
					if err != nil {
						log.Panic(err)
					}

					checkErr(err)
					return nil
				})
				checkErr(err)
			}
		}

		save := data[nodeId%3000]

		end := 0
		check := false
		if count != 0 {
			err = bc.db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(blocksBucket))
				c := b.Cursor()
				if c != nil {
					keyBytes, _ := c.Last()

					targetBytes := []byte{72, 97, 115, 104} // Hash
					for keyBytes != nil {

						if len(keyBytes) >= 4 && bytes.HasPrefix(keyBytes[:4], targetBytes) {
							var keyCheck int
							keyCheck = 0
							log.Println(keyCheck)
							for i := len(keyBytes) - 1; i >= 4; i-- {

								if keyBytes[i] == '~' {
									check = true
									keyCheck = i
								} else if !check {

									log.Println("현재 KeyBytes", string(keyBytes))
									m := math.Pow(10, float64(len(keyBytes)-i-1))
									end += int(keyBytes[i]-48) * int(m)
								}
							}

							return nil
						}
						keyBytes, _ = c.Prev() // 이전 항목으로 이동
					}
				}
				return nil
			})
			checkErr(err)

			startIndex := "Hash" + strconv.Itoa(end+1) + "~" + strconv.Itoa(end+int(3)+int(7))

			file, err := os.OpenFile("/Users/leeseonghyeon/Documents/go/go/blockchain/"+nodeID+"/"+startIndex, os.O_RDWR|os.O_CREATE, 0666)
			_, err = file.Write(save)
			checkErr(err)
			file.Close()
		} else {
			_, err := os.Stat(nodeID)
			startIndex := "Hash0~" + strconv.Itoa(9)
			//IsNotExist Folder -> mkdir Folder
			if os.IsNotExist(err) {
				err := os.Mkdir(nodeID, os.ModePerm)
				checkErr(err)
			}
			file, err := os.OpenFile("/Users/leeseonghyeon/Documents/go/go/blockchain/"+nodeID+"/"+startIndex, os.O_RDWR|os.O_CREATE, 0666)
			_, err = file.Write(save)
			checkErr(err)
			file.Close()
		}
		count++
		if count == 10 {
			break
		}
	}
	return &proto.SendResponse{}, nil
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
	start := time.Now()
	log.Println("RSEncoding")
	nodeId, err := strconv.Atoi(req.NodeId)
	checkErr(err)
	var checkList []int
	bc := NewBlockchain(req.NodeId)
	defer bc.db.Close()
	bci := bc.Iterator()

	enc, err := reedsolomon.New(7, 3) //비잔틴 장애 내성 가지도록 설계
	checkErr(err)
	//샤딩할 부분을 나눈다

	data := make([][]byte, int(10))

	for i := 0; ; i++ {
		block, err := bci.Next()
		checkErr(err)
		if block.Height/7 == int(req.Count) && block.Height%7 == 0 {

			break
		}
	}
	for i := 6; i >= 0; i-- {
		block, err := bci.Next() // 6
		//log.Println("!")
		checkErr(err)

		//log.Println("RSEncoding Block Height", block.Height)
		newBlockBytes := block.Serialize()
		// 비어있는 곳을 0으로 채운 후, newBlockBytes의 내용을 복사합니다
		data[i] = make([]byte, 2048)
		copy(data[i], newBlockBytes)
		checkList = append(checkList, block.Height)
	}

	for i := 7; i < 10; i++ {
		data[i] = make([]byte, 2048)
	}

	err = enc.Encode(data)
	checkErr(err)
	ok, err := enc.Verify(data)
	checkErr(err)
	if ok == false {
		log.Panicln("!@#!@#!@#!@#!@#!@#!@")
	}

	bci = bc.Iterator()

	for i := 0; ; i++ {
		block, err := bci.Next()
		checkErr(err)
		//log.Println(block.Height)
		if block.Height/7 == int(req.Count) && block.Height%7 == 0 {
			break
		}
	}
	save := data[nodeId%3000]
	//end := 0
	check := false
	if req.Count != 1 {
		err = bc.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			c := b.Cursor()
			if c != nil {
				keyBytes, _ := c.Last()

				targetBytes := []byte{72, 97, 115, 104}
				for keyBytes != nil {

					if len(keyBytes) >= 4 && bytes.HasPrefix(keyBytes[:4], targetBytes) {

						//log.Println(keyBytes)
						for i := len(keyBytes) - 1; i >= 4; i-- {
							//log.Println(string(keyBytes[i]))
							if keyBytes[i] == '~' {
								check = true

							} else if !check {
								//end = int(req.Count-1) * 6
							}
						}

						return nil
					}
					keyBytes, _ = c.Prev() // 이전 항목으로 이동
				}
			}
			return nil
		})
		checkErr(err)
		a := int(req.Count - 1)
		log.Println(a*7 + 6)
		startIndex := "Hash" + strconv.Itoa(a*7) + "~" + strconv.Itoa(a*7+6) // -> 이게 잘못됨
		log.Println(startIndex)
		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			b.Put([]byte(startIndex), save)

			return nil
		})
		checkErr(err)
	} else {

		startIndex := "Hash0~" + strconv.Itoa(6)
		log.Println(startIndex)
		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			b.Put([]byte(startIndex), save)

			return nil
		})
		checkErr(err)
	}
	bci = bc.Iterator()
	for i := 0; ; i++ {
		block, err := bci.Next()
		checkErr(err)
		//	log.Println(block.Height)
		if block.Height/7 == int(req.Count) && block.Height%7 == 0 {

			break
		}
	}
	for i := 0; i < 7; i++ {
		block, err := bci.Next() //6

		checkErr(err)

		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			blockInDb := b.Get(block.Hash)

			if blockInDb == nil {
				return nil
			}
			//	log.Println("삭제한 블럭 Height", DeserializeBlock(blockInDb).Height)
			err := b.Delete(block.Hash)

			if err != nil {
				log.Panic(err)
			}

			checkErr(err)
			return nil
		})
		checkErr(err)

	}
	endTime := time.Since(start)
	rs_server_process[cnt3][1] = endTime.String() // 각 노드 시간
	return &blockchain.RSEncodingResponse{}, nil
}
func (s *server) DeleteBlock(ctx context.Context, req *proto.RSEncodingRequest) (*proto.RSEncodingResponse, error) {
	var checkList []int
	bc := NewBlockchain(req.NodeId)
	defer bc.db.Close()
	bci := bc.Iterator()
	log.Println("count", req.Count)
	for i := 0; ; i++ {
		block, err := bci.Next()
		checkErr(err)
		log.Println(block.Height)
		if block.Height/7 == int(req.Count) && block.Height%7 == 0 {
			log.Println("나누어진", block.Height)
			break
		}
	}
	for i := 0; i < 7; i++ {
		block, err := bci.Next() //6
		log.Println("block:", block.Height)
		log.Println("checkList: ", checkList)
		checkErr(err)

		err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			blockInDb := b.Get(block.Hash)

			if blockInDb == nil {
				return nil
			}
			log.Println("삭제한 블럭 Height", DeserializeBlock(blockInDb).Height)
			err := b.Delete(block.Hash)

			if err != nil {
				log.Panic(err)
			}

			checkErr(err)
			return nil
		})
		checkErr(err)

	}

	return &proto.RSEncodingResponse{}, nil
}

var elapsedTime2 time.Duration
var endExcel time.Duration

func (s *server) RsReEncoding(ctx context.Context, req *proto.RsReEncodingRequest) (*proto.RsReEncodingResponse, error) {
	startTime2 := time.Now()
	nodeId, err := strconv.Atoi(req.NodeId)
	checkErr(err)

	bc := NewBlockchain(req.NodeId)
	defer bc.db.Close()

	enc, err := reedsolomon.New(7, 2) //비잔틴 장애 내성 가지도록 설계
	checkErr(err)

	var RsData [][]byte
	RsData = make([][]byte, 10)
	RsData[0] = nil
	for i := 1; i <= 7; i++ {
		RsData[i] = req.Data[i]
	}
	for i := 8; i <= 9; i++ {
		RsData[i] = make([]byte, 2048)
	}
	enc.Encode(RsData)
	var start int
	var end int
	var check int
	check = 0
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		c := b.Cursor()

		if c != nil {
			// 키 공간을 역순으로 순회하여 가장 마지막 항목을 찾음
			keyBytes, _ := c.Last()
			targetBytes := []byte{72, 97, 115, 104}
			start = 0
			end = 0

			for keyBytes != nil {
				if len(keyBytes) >= 4 && bytes.HasPrefix(keyBytes[:4], targetBytes) {

					var numberPart string
					parts := strings.Split(req.Hash[check], "Hash")
					if len(parts) > 1 {
						numberPart = parts[1] // 숫자 부분인 "0~6"을 얻음
						fmt.Println(numberPart)
					} else {
						fmt.Println("해당 형식의 문자열이 아닙니다.")
					}
					parts2 := strings.Split(numberPart, "~")
					var numbers []int
					// 분할된 각 부분에서 숫자를 int로 변환하여 배열에 추가
					for _, part := range parts2 {
						num, err := strconv.Atoi(part)
						if err != nil {
							fmt.Println("숫자로 변환할 수 없습니다:", part)
							continue
						}
						numbers = append(numbers, num)
					}

					start = numbers[0]
					end = numbers[1]
					log.Println("start", start)
					err = b.Delete([]byte("Hash" + (strconv.Itoa(start) + "~" + strconv.Itoa(end))))
					checkErr(err)
					if nodeId%3000 < 9 {
						save := RsData[nodeId%3000]
						b.Put([]byte("Hash"+(strconv.Itoa(start)+"~"+strconv.Itoa(end))), save)
					}
					check++
					return nil
				}

				keyBytes, _ = c.Prev()
				if keyBytes == nil || check == 9 {
					return nil
				}

			}
		}
		return nil
	})
	elapsedTime2 = time.Since(startTime2)
	log.Println("Rs재인코딩 끝난 시간: ", elapsedTime2)

	return &proto.RsReEncodingResponse{Success: true}, nil
}

var rs_server [100000][10]string
var cnt2 int
var rsReEncodeEndTime time.Duration

func (s *server) DataTransfer(ctx context.Context, req *proto.DataRequest) (*proto.DataResponse, error) {
	cnt2 = -1
	log.Println("DataTransfer")
	cnt := 0
	data := make([][]byte, 1000)
	enc, err := reedsolomon.New(7, 3)
	checkErr(err)
	cnt = len(req.Hash)
	var hash []string
	hash = req.Hash
	getShardTime := time.Now()
	for k := 0; k < len(knownNodes); k++ {
		serverAddress := fmt.Sprintf("localhost:%s", knownNodes[k][10:])
		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
		defer conn.Close()
		log.Println(serverAddress)
		if err != nil {
			log.Println(knownNodes[k], "에 연결 실패!")
		} else if req.NodeId != knownNodes[k][10:] { // 장애가 발생한 노드를 제외하고
			//여기서 지금 멈춤
			client := blockchain.NewBlockchainServiceClient(conn)
			cli := CLI{
				nodeID:     knownNodes[k][10:],
				blockchain: client,
			}

			request := &blockchain.GetShardRequest{
				NodeId: knownNodes[k][10:],
				Hash:   hash,
			}

			//각 노드에서 샤드 가져오기
			response, err := cli.blockchain.GetShard(context.Background(), request)
			if err != nil {
				log.Println("연결실패!", knownNodes[k])
			} else {
				bytes := response.Bytes
				list = response.List

				size := len(bytes)
				cnt = 0

				for j := 0; ; j++ {
					if cnt == size {
						break
					}
					data[cnt*10+k] = make([]byte, 2048)
					copy(data[cnt*10+k], bytes[cnt])
					cnt++
				}
			}
		}
	}
	getShardEndTime := time.Since(getShardTime)
	cnt2++
	rs_server[cnt2][0] = getShardEndTime.String()
	fmt.Printf("GetShard Total time taken: %s\n", getShardEndTime)

	var wg sync.WaitGroup

	//RsReEncoding (재인코딩) (7,2)
	rsReEncodeTime := time.Now()
	rsDecodeTime := time.Now()
	RsData := make([][]byte, 10)
	im := 0
	for j := 1; j < len(knownNodes); j++ {
		checkCnt := 0

		for i := im; i < im+10; i++ {
			if checkCnt == 0 {
				RsData[checkCnt] = nil
			} else {
				RsData[checkCnt] = data[i]
			}
			checkCnt++
		}
		im += 10
		log.Println(enc.Verify(RsData))
		enc.Reconstruct(RsData)
		log.Println(enc.Verify(RsData))
		rsDecodeEndTime := time.Since(rsDecodeTime)
		rs_server[cnt2][1] = rsDecodeEndTime.String()
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			serverAddress := fmt.Sprintf("localhost:%s", knownNodes[j][10:])
			conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Failed to connect to gRPC server: %v", err)
			}
			defer conn.Close()
			client := blockchain.NewBlockchainServiceClient(conn)
			cli := CLI{
				nodeID:     knownNodes[j][10:],
				blockchain: client,
			}
			response, err := cli.blockchain.RsReEncoding(context.Background(), &blockchain.RsReEncodingRequest{
				NodeId: knownNodes[j][10:],
				F:      2,
				NF:     7,
				Hash:   hash,
				Data:   RsData,
			})
			if err != nil {
				log.Panic(err)
			}
			if response.Success {
				log.Println(response.Success)
			}
		}(j)
	}

	wg.Wait()

	rsReEncodeEndTime = time.Since(rsReEncodeTime)
	fmt.Printf("Total time taken: %s\n", rsReEncodeEndTime)
	rs_server[cnt2][2] = rsReEncodeEndTime.String()
	rs_server_process[cnt3][3] = elapsedTime2.String()

	f := excelize.NewFile()
	sheetName := "성능평가_" + "총_시간_"
	sheet, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println(err)

	}

	// print header
	f.SetCellValue(sheetName, "A1", "GetShard Total Time")
	f.SetCellValue(sheetName, "C1", "RsDecode Total Time")
	f.SetCellValue(sheetName, "E1", "RsReEcoding Total Time")
	f.SetCellValue(sheetName, "G1", "CPU Usages")
	for i, row := range rs_server {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), row[0])
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), row[1])
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), row[2])
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), row[3])
	}
	filename := fmt.Sprintf("%s.xlsx", time.Now().Format("2006-01-02"))

	f.SetActiveSheet(sheet)
	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
	}
	fmt.Println("생성 완료")

	return &proto.DataResponse{Success: true}, nil
}

func (s *server) FindBlockTransaction(ctx context.Context, req *proto.FindChunkTransactionRequest) (*proto.FindChunkTransactionReponse, error) {
	bc := NewBlockchainRead(req.NodeId)
	defer bc.db.Close()
	bci := bc.Iterator()
	var transaction Transaction

	for {
		block, err := bci.Next()
		if err != nil {
			break
		}

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, req.VinId) == 0 {
				log.Println("Equal?")
				log.Println(tx)
				transaction = tx.TrimmedCopy()
				return &proto.FindChunkTransactionReponse{
					Transaction: convertToProtoTransaction(transaction)}, nil
			}
		}
	}

	return &proto.FindChunkTransactionReponse{}, nil

}
func (s *server) Sign(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	wallets, err := NewWallets("3000") // wallet.node_id 확인
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(req.HashId)
	tx := convertToOneTransaction(req.Transaction)
	tx2 := convertToOneTransaction(req.Tx)
	prevTXs := make(map[string]Transaction)

	prevTXs[req.HashId] = tx

	var privKey ecdsa.PrivateKey
	privKey = wallet.PrivateKey

	tx2.Sign(privKey, prevTXs)

	return &proto.ValidateResponse{Check: 1}, nil

}
func (s *server) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {

	tx := convertToOneTransaction(req.Transaction)
	tx2 := convertToOneTransaction(req.Tx)
	prevTXs := make(map[string]Transaction)
	log.Println(tx.ID)
	prevTXs[req.HashId] = tx
	log.Println(req.HashId)
	check := tx2.Verify(prevTXs)
	if check == true {
		return &proto.ValidateResponse{Check: 1}, nil
	}
	return &proto.ValidateResponse{Check: 0}, nil
}
func (s *server) FindChunkTransaction(ctx context.Context, req *proto.FindChunkTransactionRequest) (*proto.FindChunkTransactionReponse, error) {
	//log.Println("FindChunkTransaction")
	bc := NewBlockchainRead(req.NodeId)
	defer bc.db.Close()
	//log.Println("DB는 잘열렸다.")
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

					data = append(data, make([]byte, len(value)))
					data[cnt] = make([]byte, len(value))
					copy(data[cnt], value)

					cnt++

					block := DeserializeBlock(data[cnt-1])

					//	log.Println("블럭 높이:", block.Height)
					for _, tx := range block.Transactions {
						if bytes.Equal(tx.ID, req.VinId) {
							//	log.Println("TX----", tx)

							check = true

							//	log.Println(req.NodeId, "에서 발견")
							foundTx = tx.TrimmedCopy()
							return nil
						}

					}
				}
				keyBytes, _ = c.Prev() // 이전 항목으로 이동
				if keyBytes == nil {
					return nil
				}
			}

		}

		return nil
	})
	checkErr(err)
	//log.Println(check)
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
	log.Println("GetShard")
	getShardStartTime := time.Now()
	bc := NewBlockchainRead(req.NodeId)
	defer bc.db.Close()

	var data [][]byte
	var stringKey string
	stringKey = ""
	log.Println("lastKEy!!!!", []byte(stringKey))
	log.Println("Find Hash?", req.Hash)
	i := 0

	for {
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
						copiedValue := make([]byte, len(value))
						copy(copiedValue, value)
						stringKey := string(keyBytes)
						log.Println("Find Data", stringKey)
						log.Println("i:", i, "len", len(req.Hash))
						//원하는 Hash를 찾았을때
						if stringKey == req.Hash[i] {
							log.Println("Find Data & Send Data:", stringKey)
							data = append(data, copiedValue)
							i++
						}
					}
					keyBytes, _ = c.Prev()
					if keyBytes == nil || i == (len(req.Hash)) {
						return nil
					}
				}
			}
			return nil
		})
		if i == len(req.Hash) {
			break
		}
		checkErr(err)
	}
	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", os.Getpid()), "-o", "%cpu")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)

	}

	cpuUsageStr := strings.TrimSpace(string(output))
	log.Println("Current CPU Usage!!!!!!!!!!!!!!!!!!!!!!!!!!!:" + cpuUsageStr)

	getShardEndTime := time.Since(getShardStartTime)
	log.Println("getShardTime:", getShardEndTime)
	return &proto.GetShardResponse{
		Bytes: data, // 깊은 복사된 슬라이스를 반환 -> 블록 샤드 or 패리티 샤드
		List:  list,
	}, nil
}

func (s *server) GetOneShard(ctx context.Context, req *proto.GetOneShardRequest) (*proto.GetOneShardResponse, error) {

	return &proto.GetOneShardResponse{Data: nil}, nil
}
func (s *server) DeleteMempool(ctx context.Context, req *proto.DeleteMempoolRequest) (*proto.DeleteMempoolResponse, error) {
	log.Println("\n\n\n\nMEMPOLL SIZE", len(mempool))
	for key := range mempool {
		delete(mempool, key)
	}

	return &proto.DeleteMempoolResponse{}, nil
}

func (s *server) SendBlock(ctx context.Context, req *proto.SendBlockRequest) (*proto.SendBlockResponse, error) {
	//log.Println("블럭을 잘 전달받았음 ")
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

func (s *server) CheckRsEncoding(ctx context.Context, req *proto.CheckRsEncodingRequest) (*proto.CheckRsEncodingResponse, error) {
	bc := NewBlockchainRead(req.NodeId)
	defer bc.db.Close()

	var data []byte
	var lastKey string
	lastKey = ""
	log.Println("lastKEy!!!!", []byte(lastKey))

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

					data = make([]byte, 2048)
					data = value

				}
				keyBytes, _ = c.Prev() // 이전 항목으로 이동
				lastKey = string(keyBytes)
			}

		}

		return nil
	})

	checkErr(err)
	for j := 0; j < 10; j++ {
		if bytes.Equal(req.Bytes[j], data) {
			log.Println("J", j)
			return &proto.CheckRsEncodingResponse{Check: true}, nil
		}
	}
	return &proto.CheckRsEncodingResponse{Check: false}, nil

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
