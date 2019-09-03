package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/lukas1503k/msger/blockchain"
	"github.com/lukas1503k/msger/crypto"
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

func (account Account) GetPrivateKey() ecdsa.PrivateKey {
	return account.privateKey
}
func CreateAccount() *Account {
	//creates an empty account
	privKey, pubKey := CreateKeys()

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
	hash := sha256.Sum256(message)
	messageHash := hash[:]
	nounce := []byte(strconv.FormatInt(wallet.accountNounce, 10))
	messageHash = append(messageHash, nounce...)
	r, s, err := ecdsa.Sign(rand.Reader, &wallet.privateKey, messageHash)

	if err != nil {
		panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

func InitExchange(wallet *Account, toAddress []byte, transaction blockchain.Message) *blockchain.KeyExchange {
	var exchange *blockchain.KeyExchange
	if wallet.balance < getTransactionFee()*2 {
		log.Panic("Insufficient Funds")
	} else {
		curve := elliptic.P256()
		ephemeralKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			log.Panic(err)
		}
		proof := crypto.CreateProof(*wallet.privateKey, ephemeralKey)
		exchange = &blockchain.KeyExchange{wallet.accountNounce, toAddress, wallet.address, nil, wallet.publicKey, proof, false}
		exchangeBytes := blockchain.SerializeMessage(exchange)
		sig := SignTransaction(wallet, exchangeBytes)
		blockchain.SignMessage(exchangeBytes, sig)

	}

	return exchange
}

func RespExchange(wallet *Account, initialExchange *blockchain.KeyExchange) *blockchain.ExchangeResponse {
	curve := secp256k1.S256()
	ephemeralKey, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}
	proof := crypto.CreateProof(&wallet.privateKey, ephemeralKey)
	exchange := blockchain.ExchangeResponse{*initialExchange, nil, wallet.publicKey, wallet.address, proof}
	exchangeBytes := blockchain.SerializeMessage(exchange)
	sig := SignTransaction(wallet, exchangeBytes)
	blockchain.SignMessage(exchangeBytes, sig)

	return &exchange
}
