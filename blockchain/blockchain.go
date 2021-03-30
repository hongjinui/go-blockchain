package blockchain

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const (
	blocksBucket    = "block"
	dbFile          = "blockchain.db"
	genesisCoinbase = "this is genesis coinbase"
)

type Blockchain struct { // blockchain struct
	// blocks []*Block
	tip []byte
	db  *bolt.DB
}
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) { // 블록체인에 블록 추가
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("I"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("I"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func NewBlockchain(address string) *Blockchain { // 새로운 블록체인 생성
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			cbtx := NewCoinbaseTX(address, genesisCoinbase)
			genesis := NewGenesisBlock(cbtx)
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("I"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("I"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}
	return &bc

}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := BlockchainIterator{bc.tip, bc.db}
	return &bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = block.Deserialize(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash
	return block
}
func (bc *Blockchain) GetDB() *bolt.DB {

	return bc.db
}
func (bc *Blockchain) FindUTXOs(address string, amount int) (int, map[string][]int) {

	var spentTXs map[string][]int
	var unspentTXs map[string][]int

	accumulated := 0
	bci := bc.Iterator()

	for {
		block := bci.Next()

	Work:
		for _, tx := range block.Transactions {
			txid := string(tx.GetHash())

			for outid, out := range tx.Vout {
				if spentTXs[txid] != nil {
					for _, spentOut := range spentTXs[txid] {
						if spentOut == outid {
							continue
						}
					}
				}
				if out.Unlock(address) && accumulated < amount {
					accumulated += out.Value
					unspentTXs[txid] = append(unspentTXs[txid], outid)
				}
			}
		}

	}

}
