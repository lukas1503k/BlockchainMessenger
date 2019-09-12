package wallet

import (
	"crypto/aes"
	"crypto/ecdsa"
	"github.com/lukas1503k/msger/blockchain"
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

func (chain *MessageChain) AddMessageToChain(newMessage blockchain.Message, key []byte) {
	chain.messages = append(chain.messages, MessageKeyPair{newMessage, key})
	chain.messageCount += 1
}

func initChain(address []byte) MessageChain {

	chain := MessageChain{}
	chain.toAddress = address
	return chain
}

func startRatchet(chain MessageChain)

func sendMessage(account wallet.Account, chain MessageChain, message string, amount int64) *blockchain.Message {
	var ciphertext []byte

	cipher, err := aes.NewCipher(chain.currentKey)

	if err != nil {
		log.Panic(err)
	}

	cipher.Encrypt([]byte(message), ciphertext)

	return &CreateMessage(account, chain.toAddress, amount, ciphertext)

}
