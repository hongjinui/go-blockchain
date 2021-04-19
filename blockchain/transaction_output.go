package blockchain

import (
	"bytes"

	utils "github.com/hongjinui/go-blockchain/utils"
)

type TXOutput struct {
	Value        int
	ScriptPubKey []byte
}

// Lock signs the output
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	out.ScriptPubKey = pubKeyHash
}

// Unlock checks if the output can be used by the owner of the pubkey
func (out *TXOutput) Unlock(pubKeyHash []byte) bool {
	return bytes.Compare(out.ScriptPubKey, pubKeyHash) == 0
}

// NewTXOutput createa new TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}
