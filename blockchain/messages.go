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
}

type messageChain struct {
	contactID string
	messages  []message
}

func serializeMessageArray(messages []*message) [][]byte{
	var serializedMessages [][]byte

	for _, message := range messages{
		serializedMessage := serializeMessage(message)
		serializedMessages := append(serializedMessages, serializedMessage)
	}
	return serializedMessages

}

func serializeMessage(message *message) []byte {
	var msg bytes.Buffer

	msgEncoder := gob.NewEncoder(&msg)

	err := msgEncoder.Encode(message)
	fmt.Println(err)
	return msg.Bytes()

}

func deserializeMessage(data []byte) *message{

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(messages)
	fmt.Println(err)
	return &messages


}
