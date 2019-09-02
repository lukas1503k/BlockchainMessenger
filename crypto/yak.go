package YAK

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"math/big"
)

func GenerateSharedKey(k *big.Int, w *big.Int, QE ecdsa.PublicKey, W ecdsa.PublicKey) (*big.Int, *big.Int) {
	curve := elliptic.P256()
	r := new(big.Int)
	r.Add(k, w)

	gX, gY := curve.Add(QE.X, QE.Y, W.X, W.Y)

	Xab, Yab := curve.ScalarMult(gX, gY, r.Bytes())

	return Xab, Yab

}

func generateK(transaction *tx, key ecdsa.PrivateKey) *big.Int {
	txBytes := SerializeTx(transaction)
	tx := sha256.Sum256(txBytes)
	txHash := tx[:]
	txHashInt := new(big.Int)
	txHashInt.SetBytes(txHash)

	r, s := transaction.getSignature()

	K := new(big.Int)

	K.Mul(r, key.PublicKey.X)
	K.Add(K, txHashInt)
	K.Mul(K, 1/s)

	return K
}

func generateQ(k *big.Int) (*big.Int, *big.Int) {

	curve := elliptic.P256()

	Qx, Qy := curve.ScalarBaseMult(k.Bytes())

	return Qx, Qy

}

func determineExpectedK(QA ecdsa.PublicKey, QE ecdsa.PublicKey, k *big.Int) *big.Int {

	if QA.X == QE.X && QA.Y == QE.Y {
		return k
	} else {
		return k.Neg(k)
	}

}

func uncompress(Qx, Qy *big.Int, sign int) (*big.Int, *big.Int) {
	if sign < 0 {
		if signOf(Qy) < 0 {
			return Qx, Qy
		} else {
			return Qx, Qy.Neg(Qy)

		}

	} else {
		if signOf(Qy) > 0 {
			return Qx, Qy
		} else {
			return Qx, Qy.Neg(Qy)
		}
	}

}
func signOf(number *big.Int) int {
	return number.Sign()

}

func (transaction *tx) getSignature() (big.Int, big.Int) {
	r := big.Int{}
	s := big.Int{}

	sigLength := len(transaction.signature)
	r.SetBytes(transaction.signature[:sigLength])
	s.SetBytes(transaction.signature[sigLenght:])
	return r, s
}
