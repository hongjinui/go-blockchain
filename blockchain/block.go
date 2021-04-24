package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	m "github.com/hongjinui/go-blockchain/merkletree"
)

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// HashTransactions returns a hash of the transactions int the block
func (b *Block) HashTransactions() []byte {
	var transaction [][]byte

	for _, tx := range b.Transactions {
		transaction = append(transaction, tx.Serialize())

	}
	mTree := m.NewMerkleTree(transaction)
	return mTree.RootNode.Data
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block { // data, prevBlockHash를 받아 새로운 블록 생성
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	// block.SetHash()
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block { // 최초 블록 생성
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// Serialize returns a serialized Transaction
func (b Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// DeserializeBlock deserializes a block
func (b *Block) DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
