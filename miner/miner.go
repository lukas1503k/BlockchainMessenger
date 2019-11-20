package miner

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha512"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/lukas1503k/msger/blockchain"
	"math/big"
)

}



var unpackedBlocks [] interface{}

func verifyMessage(msg interface{}) bool {
	blankMessage := msg
	blankMessage.Signature = nil
	blankMessage.PublicKey = nil
	messageHash := fmt.Sprintf("%x\n", blankMessage)
	blankMessage = nil

	r := big.Int{}
	s := big.Int{}

	sigLength := len(msg.Signature())
	r.SetBytes(msg.Signature()[:sigLength])
	s.SetBytes(msg.Signature()[sigLenght:])

	msgKeyHash := sha512.Sum512(msg.PublicKey)
	msgKeyAddress := msgKeyHash[:30]
	if bytes.Compare(msgKeyAddress, msg.From) == 1 {
		return false
	}

	keyLength := len(msg.PublicKey)
	x := big.Int{}
	y := big.Int{}

	x.SetBytes(msg.publicKey[:keyLength])
	y.SetBytes(msg.publicKey[keyLength:])

	curve := secp256k1.S256()

	rawPublicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	return ecdsa.Verify(&rawPublicKey, []byte(messageHash), &r, &s)

}



func packageBlock( messages []interface{}, state WorldState) *blockchain.Blockchain {
	var verifiedMessages []interface{}
	for msg := range messages{
		if verifyMessage([msg]) == true && state.VerifyTransaction(msg)

	}




}
