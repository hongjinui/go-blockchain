package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
)

const subdity = 10

type Transaction struct {
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (tx Transaction) isCoinbaseTX() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (out *TXOutput) Unlock(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
func (in *TXInput) LockedBy(address string) bool {
	return in.ScriptSig == address

}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {

	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subdity, to}
	tx := Transaction{[]TXInput{txin}, []TXOutput{txout}}
	return &tx
}

// NewUTOXTransaction creates a new transaction
func NewUTOXTransaction(from, to string, value int, bc *Blockchain) *Transaction {

	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindUTXOs(from, value)
	if acc < value {
		log.Panic("ERROR : Not enough funds")
	}
	for txID, outs := range validOutputs {
		for _, out := range outs {
			txidbytes, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			input := TXInput{txidbytes, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{value, to})
	if acc > value {
		outputs = append(outputs, TXOutput{acc - value, to})
	}

	tx := Transaction{inputs, outputs}
	return &tx
}
