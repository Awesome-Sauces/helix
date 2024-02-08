package validium

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type Sig struct {
	Signature []byte
	Address   []byte
	ID        int
}

type Transaction struct {
	Sender     []byte
	Payload    []byte
	Originator int
	Signatures []Sig
}

type Node struct {
	ID           int
	Listener     chan Transaction
	Transactions map[string]Transaction
	Peers        map[int]chan Transaction
}

func NewNode(id int) *Node {
	return &Node{
		ID:           id,
		Listener:     make(chan Transaction, 100), // Buffer to prevent blocking
		Transactions: make(map[string]Transaction),
		Peers:        make(map[int]chan Transaction),
	}
}

func (n *Node) PeerConnect(id int, channel chan Transaction) {
	n.Peers[id] = channel
}

func (n *Node) Listen(logChan chan string) {
	for tx := range n.Listener {
		hash := crypto.Keccak256Hash(tx.Payload, tx.Sender).Hex()
		if _, exists := n.Transactions[hash]; !exists {
			n.Transactions[hash] = tx
			logChan <- fmt.Sprintf("Node %d received transaction from %d", n.ID, tx.Originator)
			for peerID, peerChan := range n.Peers {
				logChan <- fmt.Sprintf("Node %d sending to Node %d", n.ID, peerID)
				peerChan <- tx
			}
		}
	}
}

func setupNodes(nodeCount int) []*Node {
	nodes := make([]*Node, nodeCount)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}
	return nodes
}

func connectRandomPeers(nodes []*Node, peersCount int) {
	rand.Seed(time.Now().UnixNano())
	for _, node := range nodes {
		perm := rand.Perm(len(nodes))
		for _, idx := range perm {
			if len(node.Peers) >= peersCount {
				break
			}
			if nodes[idx].ID != node.ID && node.Peers[nodes[idx].ID] == nil {
				node.PeerConnect(nodes[idx].ID, nodes[idx].Listener)
			}
		}
	}
}

func TestNetworkPropagation(t *testing.T) {
	const nodeCount = 30
	const peersCount = 3 // Number of peers per node
	nodes := setupNodes(nodeCount)
	connectRandomPeers(nodes, peersCount)

	logChan := make(chan string, 10000) // Large buffer to prevent blocking
	var wg sync.WaitGroup

	for _, node := range nodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			n.Listen(logChan)
		}(node)
	}

	// Wait a bit before sending the transaction to ensure all nodes are listening
	time.Sleep(1 * time.Second)

	// Sending a transaction from a random node to test propagation
	randNode := nodes[rand.Intn(nodeCount)]
	tx := Transaction{
		Sender:     []byte("sender"),
		Payload:    []byte("payload"),
		Originator: randNode.ID,
		Signatures: []Sig{{Signature: []byte("signature"), Address: []byte("address"), ID: 1}},
	}

	randNode.Listener <- tx

	// Let the network process the transaction
	time.Sleep(5 * time.Second)
	close(logChan) // Close log channel to finish logging

	// Wait for all Listen goroutines to finish
	wg.Wait()

	// Log processing
	for logEntry := range logChan {
		fmt.Println(logEntry)
	}

	// This example demonstrates the setup and basic transaction propagation with logging.
	// Generating a text-based graph from the log entries would involve analyzing the log strings
	// to visualize the transaction flow between nodes.
}
