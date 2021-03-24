package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	utils "github.com/hongjinui/go-blockchain/utils"
)

const (
	targetBits = 15
	maxNonce   = math.MaxInt64
)

var (
	nonce int
	// isValid int
)

type ProofOfWork struct { // prooofwork  struct
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork { // proofofwork 생성 함수
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte { // mining 하기 전 데이터 직렬화

	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		utils.IntToHex(pow.block.Timestamp), // IntToHex: block에 저장되는 데이터를 16진법으로 표현하는 함수
		utils.IntToHex(int64(targetBits)),
		utils.IntToHex(int64(nonce)),
	},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) { // mining

	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data) //print formating

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 { // hashInt에 SetBytes()한 후 pow.target과 대소빇
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validation() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
