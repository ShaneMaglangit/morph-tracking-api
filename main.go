package main

import (
	"log"
	"morph-tracking-api/database"
	"morph-tracking-api/router"
	"morph-tracking-api/rpc"
	"time"
)

// BlockRange number of blocks that will be processed each iteration
const BlockRange = 500

func main() {
	db := database.New()
	rpcClient := rpc.New()

	// Run the RPC listener on the background.
	go CrawlMorph(db, rpcClient)
	router.Listen(db)
}

func CrawlMorph(db *database.AxieDB, rpcClient *rpc.Client) {
	// Create the starting and ending block for polling.
	currentBlock := db.GetLatestBlock()
	endBlock := rpc.GetLatestBlockNumber(rpcClient)

	for ; currentBlock <= endBlock; {
		// Get the logs for morphing events between the current and end block numbers.
		log.Println("Fetching block", currentBlock, "to", currentBlock+BlockRange)
		
		filter := rpc.GetMorphFilter(int64(currentBlock), int64(currentBlock+BlockRange))
		logs := rpc.GetLogs(rpcClient, filter)
		if len(logs) == 0 {
			currentBlock = GetNextStartingBlock(currentBlock, endBlock)
			endBlock = GetNextEndBlock(rpcClient, currentBlock)
			continue
		}

		// Get the timestamp of the blocks with morph event
		log.Println("Processing", len(logs), "blocks")
		blocksNumbers := rpc.GetBlocksFromLogs(logs)
		blocks := rpc.GetBlocks(rpcClient, blocksNumbers)

		// Get the morphed Axie details from the logs
		axies := rpc.GetAxieFromLogs(blocks, logs)

		// Save the results to the database
		db.SaveAxieMultiple(axies)

		currentBlock = GetNextStartingBlock(currentBlock, endBlock)
		endBlock = GetNextEndBlock(rpcClient, currentBlock)
	}
}

func GetNextStartingBlock(currentBlock uint64, endBlock uint64) uint64 {
	nextBlock := currentBlock + BlockRange
	if nextBlock > endBlock {
		return endBlock
	}
	return nextBlock
}

func GetNextEndBlock(rpcClient *rpc.Client, currentBlock uint64) uint64 {
	nextBlock := rpc.GetLatestBlockNumber(rpcClient)
	for currentBlock == nextBlock {
		time.Sleep(30 * time.Second)
		nextBlock = rpc.GetLatestBlockNumber(rpcClient)
	}
	return nextBlock
}
