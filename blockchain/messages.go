package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/lukas1503k/msger/crypto"
	"github.com/lukas1503k/msger/wallet"
	"github.com/status-im/doubleratchet"
	"math/big"
)

type Message struct {
	nounce    int64
	To        []byte
	From      []byte
	Message   []byte
	Signature []byte
	amount    float64
}

type KeyExchange struct {
	nounce     int64
	To         []byte
	From       []byte
	Signature  []byte
	PublicKey  []byte
	SchnorrZKP *crypto.SchnorrProof
	responded  bool
}

type ExchangeResponse struct {
	initialMessage KeyExchange
	Signature      []byte
	PublicKey      []byte
	Address        []byte
	SchnorrZKP     *crypto.SchnorrProof
}

func CreateMessage(account *wallet.Account, to []byte, amount float64, message doubleratchet.Message) Message {
	m := SerializeMessage(message)
	messageBlock := Message{account.AccountNounce, to, account.Address, m, nil, amount}
	messageHash := HashMessage(message)
	sig := wallet.SignTransaction(account, messageHash)

	messageBlock.Signature = sig
	return messageBlock

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

func HashMessage(message interface{}) []byte {
	messageBytes := SerializeMessage(message)
	messageHash := sha256.Sum256(messageBytes)
	return messageHash[:]
}
