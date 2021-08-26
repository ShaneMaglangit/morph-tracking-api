package rpc

import (
	"context"
	_ "embed"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"morph-tracking-api/database"
	"time"
)

type Client struct {
	requests  int
	rateLimit int
	prevStop  time.Time
	client    *ethclient.Client
}

type TxReceipt struct {
	Root              string  `json:"root"`
	Status            string  `json:"status"`
	CumulativeGasUsed string  `json:"cumulativeGasUsed"`
	LogsBloom         string  `json:"logsBloom"`
	Logs              []TxLog `json:"logs"`
	TransactionHash   string  `json:"transactionHash"`
	ContractAddress   string  `json:"contractAddress"`
	GasUsed           string  `json:"gasUsed"`
	BlockHash         string  `json:"blockHash"`
	BlockNumber       string  `json:"blockNumber"`
	TransactionIndex  string  `json:"transactionIndex"`
}

type TxLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

func New() *Client {
	client, err := ethclient.Dial("https://api.roninchain.com/rpc")
	if err != nil {
		log.Fatal(err)
	}
	return &Client{rateLimit: 1000, prevStop: time.Now(), client: client}
}

// GetClient returns the ethclient with rate limiting
func (rpc *Client) GetClient() *ethclient.Client {
	// Increment the request count everytime the client is accessed since we expect
	// a request to be done everytime this happens.
	// NOTE: DO NOT SAVE A REFERENCE TO THE CLIENT OUTSIDE
	if rpc.requests < rpc.rateLimit {
		rpc.requests++
		return rpc.client
	}
	// Sleeps the thread if we reach the limit of 1k request every 5 min
	log.Println("RPC client limit reached")
	currentTime := time.Now()
	time.Sleep(rpc.prevStop.Add(5 * time.Minute).Sub(currentTime))
	rpc.prevStop = currentTime
	rpc.requests = 0
	log.Println("RPC client resuming")
	return rpc.client
}

// GetLogs gets the logs based on the given filter.
func GetLogs(rpc *Client, filter ethereum.FilterQuery) []types.Log {
	logs, err := rpc.GetClient().FilterLogs(context.Background(), filter)
	if err != nil {
		log.Fatal("GetLogs:", err)
	}
	return logs
}

// GetBlocks gets the block details from the chain and returns the block number with its corresponding
// timestamp.
func GetBlocks(rpc *Client, blockNumbers []uint64) map[uint64]time.Time {
	// Map to keep track of block numbers recorded
	blocks := make(map[uint64]time.Time)
	for _, blockNumber := range blockNumbers {
		// Check if block number is already recorded
		if _, found := blocks[blockNumber]; !found {
			// Get the block data by its number
			block, err := rpc.GetClient().BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
			if err != nil {
				log.Fatal("GetBlocks:", err)
			}
			blocks[blockNumber] = time.Unix(int64(block.Time()), 0)
		}
	}
	return blocks
}

// GetLatestBlockNumber gets the latest block number from the chain
func GetLatestBlockNumber(rpc *Client) uint64 {
	blockNumber, err := rpc.GetClient().BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return blockNumber
}

// GetTokenIdFromEvolveLog extracts the token id from the morph event logs
func GetTokenIdFromEvolveLog(lg types.Log) uint64 {
	tokenId := big.NewInt(0)
	tokenId.SetString(lg.Topics[1].String()[2:], 16)
	return tokenId.Uint64()
}

// GetAxieFromLogs extracts the Axie morph information from the logs
func GetAxieFromLogs(blocks map[uint64]time.Time, logs []types.Log) []database.Axie {
	var axies []database.Axie
	for _, lg := range logs {
		axies = append(axies, database.Axie{
			Hash:        lg.TxHash.String(),
			BlockNumber: lg.BlockNumber,
			Timestamp:   blocks[lg.BlockNumber],
			TokenId:     GetTokenIdFromEvolveLog(lg),
		})
	}
	return axies
}

// GetMorphFilter generates the FilterQuery object to find transactions logs for successful marketplace sales
func GetMorphFilter(start int64, end int64) ethereum.FilterQuery {
	// Address of the marketplace contracts
	contractAddress := common.HexToAddress("0x32950db2a7164ae833121501c797d79e7b79d74c")
	// Topic used to determine an auction successful event
	morphTopic := common.HexToHash("0xa006fbbbc9600fe3b3757442d103355696bba0d2b8f9201852984b64d72a0a0b")
	return ethereum.FilterQuery{
		FromBlock: big.NewInt(start),
		ToBlock:   big.NewInt(end),
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{morphTopic}},
	}
}

// GetBlocksFromLogs extracts the block numbers from the logs
func GetBlocksFromLogs(logs []types.Log) []uint64 {
	var blocks []uint64
	for _, lg := range logs {
		blocks = append(blocks, lg.BlockNumber)
	}
	return blocks
}
