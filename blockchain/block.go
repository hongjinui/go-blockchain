package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct { // block struct
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func (b *Block) SetHash() { //block data 해시화

	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	header := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(header)
	b.Hash = hash[:]

}

func NewBlock(data string, prevBlockHash []byte) *Block { // data, prevBlockHash를 받아 새로운 블록 생성
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	// block.SetHash()
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

func NewGenesisBlock() *Block { // 최초 블록 생성
	return NewBlock("Genesis Block", []byte{})
}
