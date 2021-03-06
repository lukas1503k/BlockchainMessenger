package messages

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"strconv"

	"github.com/status-im/doubleratchet"
	"math/big"
)

type SchnorrProof struct {
	r big.Int
	V ecdsa.PublicKey
}

type Message struct {
	nounce    int64
	To        []byte
	From      []byte
	Message   []byte
	Signature []byte
	Publickey []byte
	Amount    float64
}

type KeyExchange struct {
	nounce     int64
	To         []byte
	From       []byte
	Signature  []byte
	PublicKey  []byte
	Amount     float64
	SchnorrZKP *SchnorrProof
	responded  bool
}

type ExchangeResponse struct {
	InitialMessage KeyExchange
	Signature      []byte
	PublicKey      []byte
	From           []byte
	SchnorrZKP     *SchnorrProof
}

func signTransaction(privKey ecdsa.PrivateKey, accountNounce int64, message []byte) []byte {
	hash := sha256.Sum256(message)
	messageHash := hash[:]
	nounce := []byte(strconv.FormatInt(accountNounce, 10))
	messageHash = append(messageHash, nounce...)
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, messageHash)

	if err != nil {
		panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

func CreateMessage(privKey ecdsa.PrivateKey, accountNounce int64, to []byte, amount float64, message doubleratchet.Message) Message {
	m := SerializeMessage(message)
	messageBlock := Message{account.AccountNounce, to, account.Address, m, nil, account.PublicKey, amount}
	messageHash := HashMessage(message)
	sig := signTransaction(privKey, accountNounce, messageHash)

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
		serializedMessage := SerializeMessage(message)
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

func DeserializeKeyExchange(data []byte) *KeyExchange {
	var decodedMessage KeyExchange
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(decodedMessage)
	fmt.Println(err)
	return &decodedMessage

}
