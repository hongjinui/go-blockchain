package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hongjinui/go-blockchain/utils"
	"golang.org/x/crypto/ripemd160"
)

const (
	version    = byte(0x00)
	walletFile = "wallet.dat"
)

//Wallet
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

type Wallets struct {
	Wallets map[string]*Wallet
}

//GetAddresses returns an array of addresses stored int the wallet file
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)

	}
	return addresses
}
func (w Wallet) GetAddress() []byte {
	public := append(w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes()...)

	publicSHA256 := sha256.Sum256(public)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	versionPayload := append([]byte{version}, publicRIPEMD160...)

	checksum := checksum(versionPayload)
	fullPayload := append(versionPayload, checksum...)
	address := utils.Base58Encode(fullPayload)
	return address

}

// SaveToFile saves the wallets to a file
func (ws Wallets) SaveToFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// NewWallet creates and returns a Wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, ecdsa.PublicKey) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}
	return *private, private.PublicKey
}

// CreateWallet adds a Wallet to Wallets
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	ws.Wallets[address] = wallet

	return address
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:4]

}

// LoadFromFile loads wallets from the file
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = wallets.Wallets
	return nil
}

// NewWallet ...
func NewWallets() *Wallets {

	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	if err != nil {
		fmt.Println("Wallets file doesn't exist")
	}
	return &wallets
}
