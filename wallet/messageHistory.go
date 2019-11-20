package wallet

import (
	"crypto/ecdsa"
	"github.com/lukas1503k/msger/blockchain"
	"github.com/lukas1503k/msger/crypto"
	"github.com/lukas1503k/msger/wallet"
	"github.com/status-im/doubleratchet"
	"log"
)

type MessageKeyPair struct {
	message blockchain.Message
	key     []byte
}

type MessageChain struct {
	toAddress      []byte
	messages       []MessageKeyPair
	messageCount   int
	ephemeralKey   ecdsa.PrivateKey
	currentKey     []byte
	ratchetSession doubleratchet.Session
}

type dhPair struct {
	privateKey doubleratchet.Key
	publicKey  doubleratchet.Key
}



func (chain *MessageChain) AddMessageToChain(newMessage blockchain.Message, key []byte) {
	chain.messages = append(chain.messages, MessageKeyPair{newMessage, key})
	chain.messageCount += 1
}

func getKey(account wallet.Account, keyExchangeInit blockchain.KeyExchange, keyExchangeResponse blockchain.ExchangeResponse, messageChain wallet.MessageChain) {

	messageChain.currentKey = crypto.GenerateSharedKey(account, keyExchangeInit, keyExchangeResponse, messageChain)
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

func convertToDHPair(key ecdsa.PrivateKey) DHPair {
	var privateKey, publicKey [32]byte
	copy(privateKey[:], key.D.Bytes())
	privateKey[0] &= 248
	privateKey[31] &= 127
	privateKey[31] |= 64

	copy(publicKey[:], key.X.Bytes())
	return dhPair{privateKey, publicKey}
}
func StartRatchet(chain MessageChain, account Account) {
	var err error
	dhPair := convertToDHPair(chain.ephemeralKey)
	chain.ratchetSession, err = doubleratchet.New(chain.toAddress, toKey(chain.currentKey), dhPair, account.Storage)
	if err != nil {
		log.Panic(err)
	}

}

func sendMessage(account wallet.Account, chain MessageChain, message string, amount float64) *blockchain.Message {

	m, err := chain.ratchetSession.RatchetEncrypt([]byte(message), nil)

	if err != nil {
		log.Panic(err)
	}

	newMessage := blockchain.CreateMessage(account, chain.toAddress, amount, m)

	return &newMessage
}
