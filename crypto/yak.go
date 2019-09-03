package YAK

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/lukas1503k/msger/blockchain"
	"github.com/lukas1503k/msger/wallet"
	"math/big"
)

func GenerateSharedKey(account wallet.Account, keyExchangeInit blockchain.KeyExchange, keyExchangeResponse blockchain.ExchangeResponse, messageChain wallet.MessageChain) (*big.Int, *big.Int) {

	var r, s, rB, w *big.Int
	var WB *ecdsa.PublicKey
	var messageHash []byte
	if bytes.Equal(keyExchangeInit.From, account.Address) {
		emptyInit := keyExchangeInit
		emptyInit.Signature = nil

		sig := keyExchangeInit.Signature
		sigLen := len(sig)
		r = sig[:sigLen]
		s = sig[sigLen:]
		rB = keyExchangeResponse.Signature[:sigLen]

		messageHash = blockchain.HashMessage(emptyInit)

		WB := keyExchangeResponse.SchnorrZKP.V

	} else {

		emptyResponse := keyExchangeResponce
		emptyResponse.Signature = nil

		sig := keyExchangeResponse.Signature
		sigLen := len(sig)
		r = sig[:sigLen]
		s = sig[sigLen:]
		rB = keyExchangeInit.Signature[:sigLen]

		messageHash = blockchain.HashMessage(emptyResponse)
		WB := keyExchangeInit.SchnorrZKP.V

	}

	w = messageChain.ephemeralKey.D

	k := GetK(r, s, messageHash, account)

	curve := secp256k1.S256()

	QA := GetQ(k, curve)
	expectedQA := uncompress(*QA, 1)

	k = determineExpectedK(*QA, expectedQA, k)

	expectedQB := GetEllipticKeyPair(rB, curve)

	QB := uncompress(*expectedQB, 1)

	c := new(big.Int)

	c.Add(k, w)

	addedX, addedY := curve.Add(QB.X, QB.Y, WB.X, WB.Y)

	finalX, finalY := curve.ScalarMult(addedX, addedY, c.Bytes())

	return finalX, finalY

}

func GetEllipticKeyPair(x *big.Int, curve elliptic.Curve) *ecdsa.PublicKey {
	y2 := new(big.Int)
	x2 := new(big.Int)
	x2.Mul(x, x)
	x3 := new(big.Int)
	x3.Mul(x2, x)

	y2.Add(x3, 7)

	y := new(big.Int)

	y.Sqrt(y2)

	return &ecdsa.PublicKey{curve, x, y}

}

func GetK(r, s *big.Int, msgHash []byte, account wallet.Account) *big.Int {

	msgInt := new(big.Int)
	msgInt.SetBytes(msgHash)

	k := new(big.Int)
	k.Mul(r, account.GetPrivateKey)
	k.Add(k, msgInt)
	k.Div(k, s)
	return k

}

func GetQ(k *big.Int, curve elliptic.Curve) *ecdsa.PublicKey {
	xCoord, yCoord := curve.ScalarBaseMult(k.Bytes())

	publicKey := ecdsa.PublicKey{Curve: curve, X: xCoord, Y: yCoord}

	return &publicKey

}

func determineExpectedK(QA ecdsa.PublicKey, QE ecdsa.PublicKey, k *big.Int) *big.Int {

	if QA.X == QE.X && QA.Y == QE.Y {
		return k
	} else {
		return k.Neg(k)
	}

}

func uncompress(Q ecdsa.PublicKey, sign int) ecdsa.PublicKey {
	//Qx,Qy are points on an elliptic curve
	Qx := Q.X
	Qy := Q.Y
	if sign < 0 {
		if signOf(Qy) < 0 {
			return ecdsa.PublicKey{curve, Qx, Qy}
		} else {
			return ecdsa.PublicKey{curve, Qx, Qy.Neg(Qy)}
		}

	} else {
		if signOf(Qy) > 0 {
			return ecdsa.PublicKey{curve, Qx, Qy}
		} else {
			return ecdsa.PublicKey{curve, Qx, Qy.Neg(Qy)}
		}
	}

}

func signOf(number *big.Int) int {
	return number.Sign()

}
