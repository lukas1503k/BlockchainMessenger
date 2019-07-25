package miner

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha512"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"math/big"
)

func verifyMessage(msg *message) bool {
	blankMessage = msg
	blankMessage.signature = nil
	blankMessage.pubKeyHash = nil
	messageHash := := fmt.Sprintf("%x\n", blankMessage)
	blankMessage = nil

	r := big.Int{}
	s := big.Int{}

	sigLength := len(msg.signature)
	r.SetBytes(msg.signature[:sigLength])
	s.SetBytes(msg.signature[sigLenght:])

	msgKeyHash := sha512.Sum512(msg.publicKey)
	msgKeyAddress := msgKeyHash[:30]
	if bytes.Compare(msgKeyAddress, msg.from) == 1 {
		return false
	}

	keyLength := len(msg.publicKey)
	x := big.Int{}
	y := big.Int{}

	x.SetBytes(msg.publicKey[:keyLength])
	y.SetBytes(msg.publicKey[keyLength:])

	curve := secp256k1.S256()

	rawPublicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	return ecdsa.Verify(&rawPublicKey, []byte(messageHash), &r, &s)

}
