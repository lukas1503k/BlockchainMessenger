package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type message struct {
	nounce  []byte
	to      string
	from    string
	message []byte
	amount  float64
}

type messageChain struct {
	contactID    string
	messages     []message
	messageCount int
}

func serializeMessageArray(messages []*message) [][]byte{
	var serializedMessages [][]byte

	for _, message := range messages{
		serializedMessage := serializeMessage(message)
		serializedMessages := append(serializedMessages, serializedMessage)
	}
	return serializedMessages

}

func serializeMessage(incomingMessage *message) []byte {
	var incMsg bytes.Buffer

	msgEncoder := gob.NewEncoder(&incMsg)

	err := msgEncoder.Encode(incomingMessage)
	fmt.Println(err)
	return incMsg.Bytes()

}

func deserializeMessage(data []byte) *message{
	var decodedMessage message
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(decodedMessage)
	fmt.Println(err)
	return &decodedMessage


}
