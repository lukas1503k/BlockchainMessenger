package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	exchangeCrypto "github.com/lukas1503k/BlockchainMessenger/crypto"
	"github.com/lukas1503k/BlockchainMessenger/messages"
	"github.com/status-im/doubleratchet"
	"log"
	"strconv"
)

type Account struct {
	PublicKey     []byte
	PrivateKey    ecdsa.PrivateKey
	Address       []byte
	AccountNounce int64
	Balance       float64
	Conversations *[]MessageChain
	Storage       doubleratchet.SessionStorage
}

func IncrementNounce(account *Account) int64 {
	nounce := account.AccountNounce
	account.AccountNounce += 1
	return nounce
}
func (account Account) GetPrivateKey() ecdsa.PrivateKey {
	return account.PrivateKey
}
func CreateAccount() *Account {
	//creates an empty account
	privKey, pubKey := CreateKeys()

	newAccount := Account{pubKey, privKey, nil, 0, 0, new([]MessageChain), new(doubleratchet.SessionStorage)}

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
	keyHash := sha512.Sum512(wallet.PublicKey)
	address := keyHash[:30]

	wallet.Address = address

}

func SignTransaction(wallet *Account, message []byte) []byte {
	hash := sha256.Sum256(message)
	messageHash := hash[:]
	nounce := []byte(strconv.FormatInt(wallet.AccountNounce, 10))
	messageHash = append(messageHash, nounce...)
	r, s, err := ecdsa.Sign(rand.Reader, &wallet.PrivateKey, messageHash)

	if err != nil {
		panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

func InitExchange(wallet *Account, toAddress []byte) *messages.KeyExchange {
	var exchange *messages.KeyExchange
	if wallet.Balance < getTransactionFee()*2 {
		log.Panic("Insufficient Funds")
	} else {
		curve := elliptic.P256()
		ephemeralKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			log.Panic(err)
		}
		proof := exchangeCrypto.CreateProof(&wallet.PrivateKey, ephemeralKey)
		exchange = &messages.KeyExchange{wallet.AccountNounce, toAddress, wallet.Address, nil, wallet.PublicKey, proof, false}
		exchangeBytes := messages.SerializeMessage(exchange)
		sig := SignTransaction(wallet, exchangeBytes)
		messages.SignMessage(exchange, sig)

	}

	return exchange
}

func RespExchange(wallet *Account, initialExchange *blockchain.KeyExchange) *blockchain.ExchangeResponse {
	curve := secp256k1.S256()
	ephemeralKey, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}
	proof := crypto.CreateProof(&wallet.PrivateKey, ephemeralKey, wallet.Address)
	exchange := messages.ExchangeResponse{*initialExchange, nil, wallet.PublicKey, wallet.Address, proof}
	exchangeBytes := messages.SerializeMessage(exchange)
	sig := SignTransaction(wallet, exchangeBytes)
	messages.SignMessage(exchangeBytes, sig)

	return &exchange
}

func getTransactionFee() float64 {
	return 0.0001
}
