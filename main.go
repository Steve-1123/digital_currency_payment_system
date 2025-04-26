package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Node struct {
	ID          string
	Addr        string
	Peers       map[string]string // PeerID -> Address
	Wallet      *Wallet
	Blockchain  *Blockchain
	Raft        *Raft
	DB          *gorm.DB
	RedisClient *redis.Client
	TxChan      chan Transaction
	mu          sync.Mutex
}

func NewNode(id, addr string) (*Node, error) {
	// Initialize MySQL
	dsn := "user:password@tcp(127.0.0.1:3306)/dcep?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Transaction{})

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Initialize components
	wallet, err := NewWallet()
	if err != nil {
		return nil, err
	}
	bc := NewBlockchain()
	raft := NewRaft(id)
	txChan := make(chan Transaction, 1000)

	return &Node{
		ID:          id,
		Addr:        addr,
		Peers:       make(map[string]string),
		Wallet:      wallet,
		Blockchain:  bc,
		Raft:        raft,
		DB:          db,
		RedisClient: redisClient,
		TxChan:      txChan,
	}, nil
}

func (n *Node) Start() {
	// Start P2P server
	listener, err := net.Listen("tcp", n.Addr)
	if err != nil {
		log.Fatalf("Failed to start node: %v", err)
	}
	defer listener.Close()

	// Start Raft
	go n.Raft.Run()

	// Process transactions
	go n.processTransactions()

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go n.handleConnection(conn)
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	// Simplified: Handle incoming transactions (in real project, use JSON decoding)
}

func (n *Node) processTransactions() {
	ctx := context.Background()
	for tx := range n.TxChan {
		if n.Raft.IsLeader() {
			// Verify and append transaction
			if n.verifyTransaction(tx) {
				n.Blockchain.AddTransaction(tx)
				n.DB.Create(&tx)
				// Update Redis cache
				n.RedisClient.Set(ctx, tx.From, n.getBalance(tx.From), 0)
				n.RedisClient.Set(ctx, tx.To, n.getBalance(tx.To), 0)
				// Broadcast to peers
				n.broadcastTransaction(tx)
			}
		}
	}
}

func (n *Node) verifyTransaction(tx Transaction) bool {
	// Verify ECDSA signature
	return n.Wallet.VerifySignature(tx.Hash(), tx.Signature, tx.From)
}

func (n *Node) broadcastTransaction(tx Transaction) {
	// Simplified: Broadcast to peers (in real project, use TCP/JSON)
}

func (n *Node) getBalance(address string) float64 {
	// Check Redis first
	ctx := context.Background()
	val, err := n.RedisClient.Get(ctx, address).Float64()
	if err == nil {
		return val
	}
	// Fallback to DB
	var balance float64
	n.DB.Model(&Transaction{}).Where("to_address = ?", address).Select("SUM(amount)").Scan(&balance)
	n.DB.Model(&Transaction{}).Where("from_address = ?", address).Select("SUM(amount)").Scan(&balance)
	return balance
}

func main() {
	node, err := NewNode("node1", ":8080")
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}
	fmt.Println("Starting DC/EP Transaction Simulator...")
	node.Start()
}
