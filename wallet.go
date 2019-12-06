package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "wallet.dat"
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	wallet := Wallet{*private, public}
	return &wallet
}

func (w Wallet) GetAddress() string {
	pubKeyHash := HashPublicKey(w.PublicKey)
	versionedPayLoad := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayLoad)
	fullPayload := append(versionedPayLoad, checksum...)
	address := hex.EncodeToString(fullPayload) 	//TODO: convert to base58 encoding
	return address
}

func HashPublicKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}

func ValidateAddress(address string) bool {
	pubKeyHash, _ := hex.DecodeString(address) //TODO: make this base 58 too
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
