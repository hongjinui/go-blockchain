package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const subdity = 50

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

//SetId sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx Transaction) isCoinbaseTX() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {

	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, []byte(data)}
	txout := TXOutput{subdity, []byte(to)}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {

	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic()
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR : Not enough funds")
	}

	//Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			input := TXInput{txID, out, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	bc.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("Transaction %x : ", tx.ID))

	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("Input %d : ", i))
		lines = append(lines, fmt.Sprintf("TXID : %x", input.Txid))
		lines = append(lines, fmt.Sprintf("Out : %d", input.Vout))
		lines = append(lines, fmt.Sprintf("Script : %x ", input.ScriptSig))

	}
	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("Output %d : ", i))
		lines = append(lines, fmt.Sprintf("Value : %d", output.Value))
		lines = append(lines, fmt.Sprintf("Script : %x", output.ScriptPubKey))

	}

	return strings.Join(lines, "\n")

}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, []byte{}})
	}
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.ScriptPubKey})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.isCoinbaseTX() {
		return
	}
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR : Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].ScriptSig = prevTx.Vout[vin.Vout].ScriptPubKey
		txCopy.SetID()
		txCopy.Vin[inID].ScriptSig = []byte{}
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signatur := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].ScriptSig = append(signatur, tx.Vin[inID].ScriptSig...)
	}
}

// Verify verifies signatures of Transaction inputs
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	sigLen := 64

	if tx.isCoinbaseTX() {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR : Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].ScriptSig = prevTx.Vout[vin.Vout].ScriptPubKey
		txCopy.SetID()
		txCopy.Vin[inID].ScriptSig = []byte{}

		signature := vin.ScriptSig[:sigLen]
		pubKey := vin.ScriptSig[sigLen:]

		r := big.Int{}
		s := big.Int{}
		r.SetBytes(signature[:(sigLen / 2)])
		s.SetBytes(signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(pubKey)
		x.SetBytes(pubKey[:(keyLen / 2)])
		y.SetBytes(pubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}

	}
	return true
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}
