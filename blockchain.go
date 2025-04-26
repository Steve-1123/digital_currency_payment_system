package main

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

type Transaction struct {
	ID        uint   `gorm:"primaryKey"`
	From      string `gorm:"index"`
	To        string `gorm:"index"`
	Amount    float64
	Timestamp time.Time
	Signature []byte
}

func (tx *Transaction) Hash() []byte {
	data := fmt.Sprintf("%s%s%f%s", tx.From, tx.To, tx.Amount, tx.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

type Block struct {
	Index     int
	Timestamp time.Time
	PrevHash  string
	Hash      string
	Txs       []Transaction
}

type Blockchain struct {
	Blocks []*Block
	mu     sync.Mutex
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{createGenesisBlock()},
	}
}

func createGenesisBlock() *Block {
	return &Block{
		Index:     0,
		Timestamp: time.Now(),
		PrevHash:  "0",
		Hash:      "genesis",
		Txs:       []Transaction{},
	}
}

func (bc *Blockchain) AddTransaction(tx Transaction) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	// Simplified: Add to pending transactions (in real project, create new block)
}
