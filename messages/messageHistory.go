package messages

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/lukas1503k/BlockchainMessenger/messages"
	"github.com/lukas1503k/BlockchainMessenger/wallet"
	"github.com/status-im/doubleratchet"
	"log"
	"math/big"
)

type MessageKeyPair struct {
	message Message
	key     []byte
}

type MessageChain struct {
	toAddress      []byte
	messages       []MessageKeyPair
	messageCount   int
	EphemeralKey   ecdsa.PrivateKey
	currentKey     []byte
	ratchetSession doubleratchet.Session
}

type dhPair struct {
	privateKey doubleratchet.Key
	publicKey  doubleratchet.Key
}

func GenerateSharedKey(account wallet.Account, keyExchangeInit KeyExchange, keyExchangeResponse ExchangeResponse, messageChain MessageChain) []byte {

	var r, s, rB, w *big.Int
	var WB *ecdsa.PublicKey
	var messageHash []byte
	if bytes.Equal(keyExchangeInit.From, account.Address) {
		emptyInit := keyExchangeInit
		emptyInit.Signature = nil

		sig := keyExchangeInit.Signature
		sigLen := len(sig)
		r.SetBytes(sig[:sigLen])
		s.SetBytes(sig[sigLen:])
		rB.SetBytes(keyExchangeResponse.Signature[:sigLen])

		messageHash = messages.HashMessage(emptyInit)

		WB = &keyExchangeResponse.SchnorrZKP.V

	} else {

		emptyResponse := keyExchangeResponse
		emptyResponse.Signature = nil

		sig := keyExchangeResponse.Signature
		sigLen := len(sig)
		r.SetBytes(sig[:sigLen])
		s.SetBytes(sig[sigLen:])
		rB.SetBytes(keyExchangeInit.Signature[:sigLen])

		messageHash = messages.HashMessage(emptyResponse)
		WB = &keyExchangeInit.SchnorrZKP.V

	}

	w = messageChain.EphemeralKey.D

	k := GetK(r, s, messageHash, account.GetPrivateKey().X)

	curve := secp256k1.S256()

	QA := GetQ(k, curve)
	expectedQA := uncompress(*QA, 1)

	k = determineExpectedK(*QA, expectedQA, k)

	expectedQB := GetEllipticKeyPair(rB, curve)

	QB := uncompress(*expectedQB, 1)

	c := new(big.Int)

	c.Add(k, w)

	addedX, addedY := curve.Add(QB.X, QB.Y, WB.X, WB.Y)

	finalX, finalY := curve.ScalarMult(addedX, addedY, c.Bytes())

	finalBytes32 := append(finalX.Bytes(), finalY.Bytes()...)
	finalBytes := finalBytes32[:]

	return finalBytes
}

func (chain *MessageChain) AddMessageToChain(newMessage Message, key []byte) {
	chain.messages = append(chain.messages, MessageKeyPair{newMessage, key})
	chain.messageCount += 1
}

func getKey(account wallet.Account, keyExchangeInit KeyExchange, keyExchangeResponse ExchangeResponse, messageChain MessageChain) {

	messageChain.currentKey = GenerateSharedKey(account, keyExchangeInit, keyExchangeResponse, messageChain)
}

func initChain(address []byte) MessageChain {

	chain := MessageChain{}
	chain.toAddress = address
	return chain
}

func toKey(keyBytes []byte) doubleratchet.Key {
	var key doubleratchet.Key
	copy(key[:], keyBytes)

	return key
}

func convertToDHPair(key ecdsa.PrivateKey) dhPair {
	var privateKey, publicKey [32]byte
	copy(privateKey[:], key.D.Bytes())
	privateKey[0] &= 248
	privateKey[31] &= 127
	privateKey[31] |= 64

	copy(publicKey[:], key.X.Bytes())
	return dhPair{privateKey, publicKey}
}
func StartRatchet(chain MessageChain, account wallet.Account) {
	var err error
	dhPair := convertToDHPair(chain.EphemeralKey)
	chain.ratchetSession, err = doubleratchet.New(chain.toAddress, toKey(chain.currentKey), dhPair, account.Storage)
	if err != nil {
		log.Panic(err)
	}

}

func sendMessage(account wallet.Account, chain MessageChain, message string, amount float64) messages.Message {

	m, err := chain.ratchetSession.RatchetEncrypt([]byte(message), nil)

	if err != nil {
		log.Panic(err)
	}

	newMessage := messages.CreateMessage(account.PrivateKey, account.AccountNounce, chain.toAddress, amount, m)

	return newMessage
}
