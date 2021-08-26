package database

import "time"

type Axie struct {
	Hash        string    `json:"hash,omitempty" db:"hash"`
	BlockNumber uint64    `json:"blockNumber,omitempty" db:"blockNumber"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	TokenId     uint64    `json:"id" db:"tokenId"`
}
