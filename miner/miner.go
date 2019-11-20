package miner

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha512"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/lukas1503k/msger/blockchain"
	"github.com/lukas1503k/msger/blockchain/crypto"
	"math/big"
	"reflect"
)


type messageInterface interface {
	InitialMessage() blockchain.KeyExchange
	Signature()      []byte
	PublicKey()      []byte
	From()        []byte
	To()		[]byte
	Nounce()	uint32
	Amount()	float64
	SchnorrZKP() crypto.SchnorrProof
}
func (state WorldState)verifyMessage(msg messageInterface) bool {
	blankMessage := msg
	blankMessage.Signature() = nil
	messageHash := blockchain.HashMessage(blankMessage)

	r := big.Int{}
	s := big.Int{}

	sigLength := len(msg.Signature())
	r.SetBytes(msg.Signature()[:sigLength])
	s.SetBytes(msg.Signature()[sigLength:])

	msgKeyHash := sha512.Sum512(msg.PublicKey())
	msgKeyAddress := msgKeyHash[:30]
	if bytes.Compare(msgKeyAddress, msg.From()) != 0 {
		return false
	}

	keyLength := len(msg.PublicKey())
	x := big.Int{}
	y := big.Int{}

	x.SetBytes(msg.PublicKey()[:keyLength])
	y.SetBytes(msg.PublicKey()[keyLength:])

	curve := secp256k1.S256()

	rawPublicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	address := &msg.InitialMessage()
	if !reflect.ValueOf(msg.InitialMessage()).IsNil(){
		return ecdsa.Verify(&rawPublicKey, []byte(messageHash), &r, &s) && state.VerifyAccountState(msg.From(), msg.Nounce(), msg.Amount())
	}else{
		keyExchange, err := state.chain.GetMessage(msg.InitialMessage)

		if err != nil{
			return false
		}

		if bytes.Compare(keyExchange.To,msg.From) != 0 || bytes.Compare(keyExchange.From, msg.To)  != 0|| keyExchange.Responded == true || msg.SchnorrZKP() == nil{
			return false
		}

		return ecdsa.Verify(&rawPublicKey, []byte(exchangeHash), &r, &s)
	}
}

func (state WorldState) MineBlock(){
	verfiedMessages := make(chan messageInterface)
	for len(verfiedMessages) < 3 {
		pass
	}

	newBlock := blockchain.CreateBlock(state.Chain.currentHash, verfiedMessages, state.Chain.length)
	blockchain.AddBlockToChain(state.Chain, newBlock)





}