package blockchain

import "log"

const subdity = 10

type Transaction struct {
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	Txid      int
	Vout      int
	ScriptSig string
}
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (tx Transaction) isCoinbaseTX() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Txid == -1 && tx.Vin[0].Vout == -1
}

func (out *TXOutput) Unlock(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func NewCoinbaseTX(to, data string) *Transaction {

	if data == "" {
		data = "Coinbase"
	}
	txin := TXInput{-1, -1, data}
	txout := TXOutput{subdity, to}
	tx := Transaction{[]TXInput{txin}, []TXOutput{txout}}
	return &tx
}

func NewUTOXTransaction(from, to string, value int) *Transaction {

	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := s.findUnspentOutputs(from, value)
	if acc < value {
		log.Panic("ERROR : Not enough funds")
	}
	for txid, outs := range validOutputs {
		for _, out := range outs {
			input := TXInput{txid, out, from}
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
