package blockchain

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)

//Wallet
type Wallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

func (w Wallet) GetAddress() []byte {
	publicSHA256 := sha256.Sum256(w.PublicKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	versionPayload := append([]byte{version}, publicRIPEMD160...)

	checksum := checksum(versionPayload)
	fullPayload := append(versionPayload, checksum...)
	address := Base58Encode(fullPayload)
	return address

}

// NewWallet ...
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() ([]byte, []byte) {

	curve := elliptic.P256()
	private, x, y, err := elliptic.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}
	public := append(x.Bytes(), y.Bytes()...)

	return private, public
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:4]

}
