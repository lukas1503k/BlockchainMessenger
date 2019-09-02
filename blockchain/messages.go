package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/lukas1503k/msger/crypto"
	"math/big"
)

type Message struct {
	nounce    int64
	to        string
	from      string
	message   []byte
	signature []byte
	amount    float64
}

type KeyExchange struct {
	nounce     int64
	to         []byte
	from       []byte
	signature  []byte
	publicKey  []byte
	SchnorrZKP *crypto.SchnorrProof
	responded  bool
}

type ExchangeResponse struct {
	initialMessage KeyExchange
	Signature      []byte
	publicKey      []byte
	address        []byte
	SchnorrZKP     *crypto.SchnorrProof
}

func SignMessage(msg *Message, signature []byte) {
	msg.Signature = signature
}

func GetRSValues(signature []byte) (big.Int, big.Int) {
	sigLen := len(signature)
	rBytes := signature[:sigLen]
	sBytes := signature[sigLen:]

	var r, s big.Int

	r.SetBytes(rBytes)
	s.SetBytes(sBytes)

	return r, s
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
