package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"os"

	"github.com/sleeg00/blockchain/proto"
	blockchain "github.com/sleeg00/blockchain/proto"
	"google.golang.org/grpc"
)

// CLI responsible for processing command line arguments
type CLI struct {
	nodeID     string
	blockchain blockchain.BlockchainServiceClient
}

var failNodes = []string{}
var failNodesCheck int

var NF int
var f int
var k int
var p int

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
	fmt.Println("  node -node")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

var nodeID string

func main() {

	cli := CLI{}

	fmt.Println("Enter the new value of nodeID:")
	fmt.Scanln(&nodeID)

	cli.Run()

}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {

	fmt.Printf("New nodeID: %s\n", nodeID)
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	nodeCmd := flag.NewFlagSet("node", flag.ExitOnError)
	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	cpBlockCmd := flag.NewFlagSet("cp", flag.ExitOnError)
	sendNumberOfBlockCmd := flag.NewFlagSet("block", flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")
	node := nodeCmd.String("node", "", "Node is Setting")
	sendNumberOfBlock := sendNumberOfBlockCmd.Int("block", 0, "OK")

	cpBlockCmd.String("cp", "", "ok")

	switch os.Args[1] {
	case "cp":
		err := cpBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "block":
		err := sendNumberOfBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "node":
		err := nodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "file":
		err := fileCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if cpBlockCmd.Parsed() {
		for i := 1; i < len(knownNodes); i++ {

			sourceFile, err := os.Open("blockchain_3000.db")
			checkErr(err)
			defer sourceFile.Close()

			destFile, err := os.Create("blockchain_" + knownNodes[i][10:] + ".db")
			if err != nil {
				log.Println(err)
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, sourceFile)
			checkErr(err)
		}
		log.Println("복사완료")
	}
	if sendNumberOfBlockCmd.Parsed() {
		for i := 1; i < len(knownNodes); i++ {

			sourceFile, err := os.Open("blockchain_3000.db")
			checkErr(err)
			defer sourceFile.Close()

			destFile, err := os.Create("blockchain_" + knownNodes[i][10:] + ".db")
			if err != nil {
				log.Println(err)
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, sourceFile)
			checkErr(err)
		}

		startTime := time.Now()

		for k = 0; k < *sendNumberOfBlock; k++ {
			log.Println("k", k)

			for i := 0; i < 10; i++ {

				nodeID = "3000"
				originalArgs := os.Args
				os.Args = []string{
					originalArgs[0],
					"send",
					"-from",
					"1Jw6XpfDDAZ3UK37VRytkQvd8VcVM2HC6L",
					"-to",
					"1J6GzxwAfgRdZ4yrqSbos96J3Z6mjmCpRP",
					"-amount",
					"1",
					"-mine",
				}

				cli.Run()
				os.Args = originalArgs
			}

		}

		elapsedTime := time.Since(startTime)
		fmt.Printf("Total time taken: %s\n", elapsedTime)
	}
	if fileCmd.Parsed() {
		for i := 0; i < len(knownNodes); i++ {
			serverAddress := fmt.Sprintf("localhost:%s", knownNodes[i][10:])

			conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

			if err != nil {
				log.Println(knownNodes[i], "에 연결 실패!")
				failNodes = append(failNodes, knownNodes[i][10:])
				failNodesCheck++

			} else {
				client := blockchain.NewBlockchainServiceClient(conn)
				cli := CLI{
					nodeID:     nodeID,
					blockchain: client,
				}

				request := &proto.SendRequest{
					NodeTo: knownNodes[i][10:],
				}
				response, err := cli.blockchain.SaveFileSystem(context.Background(), request)
				if response != nil {

				}

				if err != nil {
					log.Println(knownNodes[i], "에 연결 실패!")
					failNodes = append(failNodes, knownNodes[i][10:])
					failNodesCheck++

				}
			}
			conn.Close()
		}
	}
	if nodeCmd.Parsed() {
		nodeID = *node
		serverAddress := fmt.Sprintf("localhost:%s", nodeID)

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

		cli.startNode(cli.nodeID, "")

	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress, nodeID)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		if nodeID == "" {
			fmt.Printf("NODE_ID env. var is not set!")
			os.Exit(1)
		}
		serverAddress := fmt.Sprintf("localhost:%s", nodeID)

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

		cli.createBlockchain(*createBlockchainAddress, nodeID)
	}

	if createWalletCmd.Parsed() {
		if nodeID == "" {
			fmt.Printf("NODE_ID env. var is not set!")
			os.Exit(1)
		}
		serverAddress := fmt.Sprintf("localhost:%s", nodeID)

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
		cli.createWallet(nodeID)
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses(nodeID)
	}

	if printChainCmd.Parsed() {
		printChain(nodeID)
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO(nodeID)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		if *sendMine {

			failNodesCheck = 0
			newblock, bc := send(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)

			invaildNodeCount := 10 - failNodesCheck
			f = (invaildNodeCount - 1) / 3
			NF = invaildNodeCount - f
			log.Println(newblock.Height)

			log.Println("sendBlock")
			//만약 블럭이 7개가 쌓였으면 인코딩한다

			// 새로운 슬라이스를 만들고 txs의 값을 복사
			var protoTransactions []*proto.Transaction
			protoTransactions = makeClientTransactions(newblock.Transactions)
			var byte []byte

			for i := 0; i < len(knownNodes); i++ {

				serverAddress := fmt.Sprintf("localhost:%s", knownNodes[i][10:])

				conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

				if err != nil {
					log.Println(knownNodes[i], "에 연결 실패!")
					failNodes = append(failNodes, knownNodes[i][10:])
					failNodesCheck++

				} else {
					client := blockchain.NewBlockchainServiceClient(conn)
					cli := CLI{
						nodeID:     nodeID,
						blockchain: client,
					}

					block := &proto.Block{
						Timestamp:     newblock.Timestamp,
						Transactions:  protoTransactions,
						PrevBlockHash: newblock.PrevBlockHash,
						Hash:          newblock.Hash,
						Nonce:         int32(newblock.Nonce),
						Height:        int32(newblock.Height),
					}
					request := &proto.ConnectServerRequest{
						NodeId: knownNodes[i][10:],
						Block:  block,
					}
					response, err := cli.blockchain.Connect(context.Background(), request)
					byte = response.Byte
					if response != nil {
						newblock.PrevBlockHash = byte
					}

					bc.db.Close()

					// 서버에 보낼 요청 메시지 생성
					requestSend := &proto.SendRequest{
						NodeTo: knownNodes[i][10:],
					}

					responseSend, err := cli.blockchain.Send(context.Background(), requestSend)

					if err != nil {
						log.Println(knownNodes[i], "에 연결 실패!")
						failNodes = append(failNodes, knownNodes[i][10:])
						failNodesCheck++
						log.Println(responseSend)
					}
					conn.Close()
				}

			}

			if newblock.Height%7 == 0 && newblock.Height != 0 {
				log.Println("RsEncoding!!!!!")
				RsEncoding(int32(newblock.Height/7), int32(3), int32(7))
			}
			// Your existing code...

		} else {
			failNodesCheck = 0
			newblock, bc := send(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)

			invaildNodeCount := 10 - failNodesCheck
			f = (invaildNodeCount - 1) / 3
			NF = invaildNodeCount - f
			log.Println(newblock.Height)

			log.Println("sendBlock")
			//만약 블럭이 7개가 쌓였으면 인코딩한다

			// 새로운 슬라이스를 만들고 txs의 값을 복사
			var protoTransactions []*proto.Transaction
			protoTransactions = makeClientTransactions(newblock.Transactions)
			var byte []byte

			for i := 0; i < len(knownNodes); i++ {

				serverAddress := fmt.Sprintf("localhost:%s", knownNodes[i][10:])

				conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

				if err != nil {
					log.Println(knownNodes[i], "에 연결 실패!")
					failNodes = append(failNodes, knownNodes[i][10:])
					failNodesCheck++

				} else {
					client := blockchain.NewBlockchainServiceClient(conn)
					cli := CLI{
						nodeID:     nodeID,
						blockchain: client,
					}

					block := &proto.Block{
						Timestamp:     newblock.Timestamp,
						Transactions:  protoTransactions,
						PrevBlockHash: newblock.PrevBlockHash,
						Hash:          newblock.Hash,
						Nonce:         int32(newblock.Nonce),
						Height:        int32(newblock.Height),
					}
					request := &proto.ConnectServerRequest{
						NodeId: knownNodes[i][10:],
						Block:  block,
					}
					response, err := cli.blockchain.Connect(context.Background(), request)
					byte = response.Byte
					if byte != nil {
						newblock.PrevBlockHash = byte
					}

					bc.db.Close()
					// 서버에 보낼 요청 메시지 생성
					requestSend := &proto.SendRequest{
						NodeTo: knownNodes[i][10:],
					}

					responseSend, err := cli.blockchain.Send(context.Background(), requestSend)

					if err != nil {
						log.Println(knownNodes[i], "에 연결 실패!")
						failNodes = append(failNodes, knownNodes[i][10:])
						failNodesCheck++
						log.Println(responseSend)
					}
				}
				conn.Close()
			}

			/*
				//-----모든 노드 mempool에 TX를 저장시킨다. -> UTXO도 업데이트 했다. //Block시도 확인해야한다.
				tx := sendTrsaction(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)

				protoTx := makeOneTransaction(tx)

				for i := 0; i < len(knownNodes); i++ {

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

					cli.requestTransaction(*sendFrom, *sendTo, *sendAmount, knownNodes[i][10:], *sendMine, protoTx, cli.nodeID)

				}

				//-----모든 노드 mempool에 TX를 저장시킨다.
			*/
		}
	}
	if startNodeCmd.Parsed() {
		if nodeID == "" {
			fmt.Printf("NODE_ID env. var is not set!")
			os.Exit(1)
		}
		serverAddress := fmt.Sprintf("localhost:%s", nodeID)

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
		cli.startNode(cli.nodeID, *startNodeMiner)
	}
}

// 12aTcP7x7PxZcqs7DbsPUS1NY8HZcaVwqV
// 1K6BBBMDJVEjP4ZdBMNvN2jKVc2CeHTEWA
func convertFromProtoTransaction(ptx *blockchain.Transaction) Transaction {
	return Transaction{
		ID:   ptx.Id,
		Vin:  convertFromProtoInputs(ptx.Vin),
		Vout: convertFromProtoOutputs(ptx.Vout),
	}
}

func convertFromProtoInputs(protoInputs []*blockchain.TXInput) []TXInput {
	var inputs []TXInput
	for _, input := range protoInputs {
		inputs = append(inputs, TXInput{
			Txid:      input.Txid,
			Vout:      int(input.Vout),
			Signature: input.Signature,
			PubKey:    input.PubKey,
		})
	}
	return inputs
}

func convertFromProtoOutputs(protoOutputs []*blockchain.TXOutput) []TXOutput {
	var outputs []TXOutput
	for _, output := range protoOutputs {
		outputs = append(outputs, TXOutput{
			Value:      int(output.Value),
			PubKeyHash: output.PubKeyHash,
		})
	}
	return outputs
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}
