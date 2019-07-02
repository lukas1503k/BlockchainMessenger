package blockchain

import (
	"bytes"
	"encoding/gob"
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

func Serialize(messages []*message) []byte {
	var msg bytes.Buffer

	msgEncoder := gob.NewEncoder(&msg)

	msgEncoder.Encode(messages)

	return msg.Bytes()

}
