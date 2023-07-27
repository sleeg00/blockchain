package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

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

func main() {

	cli := CLI{}

	cli.Run()
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	nodeID := "3000"

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

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")
	node := nodeCmd.String("node", "", "Node is Setting")
	switch os.Args[1] {
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
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if nodeCmd.Parsed() {
		err := os.Setenv("nodeID", *node)
		if err != nil {
			log.Println("노드 번호 설정 실패")
			return
		}
		log.Println("node is setting : ", *node)
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
		cli.printChain(nodeID)
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
			newBlockPtr := send(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)

			// 새로운 슬라이스를 만들고 txs의 값을 복사
			var protoTransactions []*proto.Transaction

			for _, tx := range newBlockPtr.Transactions {

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

			for i := 0; i < len(knownNodes); i++ {
				var newblock BlockCopy
				newblock.Hash = newBlockPtr.Hash
				newblock.PrevBlockHash = newBlockPtr.PrevBlockHash
				newblock.Transactions = make([]*Transaction, len(newBlockPtr.Transactions))
				copy(newblock.Transactions, newBlockPtr.Transactions)
				newblock.Height = newBlockPtr.Height
				newblock.Nonce = newBlockPtr.Nonce
				log.Println(knownNodes[i])

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

				block := &proto.Block{
					Timestamp:     newblock.Timestamp,
					Transactions:  protoTransactions,
					PrevBlockHash: newblock.PrevBlockHash,
					Hash:          newblock.Hash,
					Nonce:         int32(newblock.Nonce),
					Height:        int32(newblock.Height),
				}

				cli.request(*sendFrom, *sendTo, *sendAmount, knownNodes[i][10:], *sendMine, block, cli.nodeID)

			}

		} else {

			//-----모든 노드 mempool에 TX를 저장시킨다. -> UTXO도 업데이트 했다. //Block시도 확인해야한다.
			tx := sendTrsaction(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)
			log.Println(tx)
			var protoTransactions *proto.Transaction

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

			var wg sync.WaitGroup //고루틴이 완료되기를 기다리기 위한 준비를 합니다.
			for i := 0; i < len(knownNodes); i++ {
				wg.Add(1) // 고루틴의 수를 증가시킵니다.
				go func(index int) {
					defer wg.Done() // 해당 고루틴이 끝나면 WaitGroup에서 하나를 차감합니다.

					serverAddress := fmt.Sprintf("localhost:%s", knownNodes[index][10:])

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

					cli.requestTransaction(*sendFrom, *sendTo, *sendAmount, knownNodes[index][10:], *sendMine, protoTransactions, cli.nodeID)

				}(i)

			}
			wg.Wait()
			//-----모든 노드 mempool에 TX를 저장시킨다.
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

func convertFromProtoTransaction(ptx *blockchain.Transaction) *Transaction {
	return &Transaction{
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

type BlockCopy struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
}
