package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/lukas1503k/msger/wallet"
)

type Message struct {
	nounce  int64
	to      string
	from    string
	message []byte
	amount  float64
}

type KeyExchange struct {
	nounce    int64
	to        []byte
	from      []byte
	signature []byte
	publicKey []byte
}

type ExchangeResponse struct {
	initialMessage KeyExchange
	signature      []byte
	publicKey      []byte
}

func InitiateExchange(initAccount *wallet.Account, destination []byte) *KeyExchange {
	msg := new(KeyExchange)
	msg.from = initAccount.GetAddress()
	msg.nounce = initAccount.accountNounce
	initAccount.accountNounce += 1
	msg.to = destination
	msg.publicKey = initAccount.publicKey
	msg.signature = wallet.SignTransaction(initAccount, SerializeMessage(msg))

	return msg

}

func SerializeMessageArray(messages []*Message) [][]byte {
	var serializedMessages [][]byte

	for _, message := range messages {
		serializedMessage := serializeMessage(message)
		serializedMessages = append(serializedMessages, serializedMessage)
	}
	return serializedMessages

}

func SerializeMessage(incomingMessage interface{}) []byte {
	var incMsg bytes.Buffer

	msgEncoder := gob.NewEncoder(&incMsg)

	err := msgEncoder.Encode(incomingMessage)
	fmt.Println(err)
	return incMsg.Bytes()

}

func DeserializeMessage(data []byte) *Message {
	var decodedMessage Message
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(decodedMessage)
	fmt.Println(err)
	return &decodedMessage

}
