package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"log"
	"strconv"
)

type Account struct {
	publicKey     []byte
	privateKey    ecdsa.PrivateKey
	address       []byte
	accountNounce int64
	balance       float64
}

func CreateAccount() *Account {
	//creates an empty account
	privKey, pubKey := getKeys()

	newAccount := Account{pubKey, privKey, nil, 0, 0}

	return &newAccount
}

func CreateKeys() (ecdsa.PrivateKey, []byte) {
	// creates a unique key pair for the wallet
	curve := secp256k1.S256()
	newPrivateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	newPublicKey := append(newPrivateKey.PublicKey.X.Bytes(), newPrivateKey.PublicKey.Y.Bytes()...)

	return *newPrivateKey, newPublicKey

}

func GetAddress(wallet *Account) {
	keyHash := sha512.Sum512(wallet.publicKey)
	address := keyHash[:30]

	wallet.address = address

}

func SignTransaction(wallet *Account, message []byte) []byte {
	messageHash := sha256.Sum256(message)
	nounce := []byte(strconv.FormatInt(wallet.accountNounce, 10))
	messageHash = append(messageHash, nounce...)
	r, s, err := ecdsa.Sign(rand.Reader, wallet.privateKey, messageHash)

	if err != nil {
		panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

func InitConversation(wallet *Account, address []byte) {
	if wallet.balance < getTransactionFee()*2 {
		log.Panic("InsuffienctFunds")
	} else {

	}

}
