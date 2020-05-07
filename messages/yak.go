package messages

import (
	"crypto/ecdsa"
	"crypto/elliptic"

	"math/big"
)

func GetEllipticKeyPair(x *big.Int, curve elliptic.Curve) *ecdsa.PublicKey {
	y2 := new(big.Int)
	x2 := new(big.Int)
	x2.Mul(x, x)
	x3 := new(big.Int)
	x3.Mul(x2, x)

	temp := new(big.Int)
	temp.Set(7)
	y2.Add(x3, temp)

	y := new(big.Int)

	y.Sqrt(y2)

	return &ecdsa.PublicKey{curve, x, y}

}

func GetK(r, s *big.Int, msgHash []byte, privKey *big.Int) *big.Int {

	msgInt := new(big.Int)
	msgInt.SetBytes(msgHash)

	k := new(big.Int)
	k.Mul(r, privKey)
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
