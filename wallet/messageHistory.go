package wallet

import (
	"crypto/ecdsa"
	"github.com/lukas1503k/msger/blockchain"
)

type MessageKeyPair struct {
	message blockchain.Message
	key     []byte
}

type MessageChain struct {
	contactID    []byte
	messages     []MessageKeyPair
	messageCount int
	ephemeralKey ecdsa.PrivateKey
}

func (chain *MessageChain) AddMessageToChain(newMessage blockchain.Message, key []byte) {
	chain.messages = append(chain.messages, MessageKeyPair{newMessage, key})
	chain.messageCount += 1
}
