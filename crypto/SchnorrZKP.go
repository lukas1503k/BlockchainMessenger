package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"math/big"
)

type SchnorrProof struct {
	r big.Int
	V ecdsa.PublicKey
}

//for this proof we assume both parties are using the same curve which has been decided ahead of time
var curve = secp256k1.S256()

func createChallenge(a *ecdsa.PrivateKey, v *ecdsa.PrivateKey, userID []byte) ([]byte, *big.Int) {

	//gets public keys from private keys
	V := v.PublicKey.X.Bytes()
	A := a.PublicKey.X.Bytes()

	//gets curve order
	parameters := curve.Params()
	n := parameters.N

	generatorPointX := parameters.Gx
	generatorPointXBytes := generatorPointX.Bytes()

	challengeBytes := append(generatorPointXBytes, V...)
	challengeBytes = append(challengeBytes, A...)

	challenge := sha256.Sum256(challengeBytes)
	c := challenge[:]
	return c, n
}

func CreateProof(a *ecdsa.PrivateKey, v *ecdsa.PrivateKey, userID []byte) *SchnorrProof {
	/*
		Takes in two private keys and returns a Schnorr ZKP
	*/

	//calls the helper function CreateChallenge
	c, n := createChallenge(a, v, userID)

	//gets public keys
	V := v.PublicKey

	//creates a numerical representation of the challenge
	cN := new(big.Int)
	cN.SetBytes(c)

	//operations calculating r = v - a*c mod n
	r := new(big.Int)
	r.Mul(a.D, cN)
	r.Sub(v.D, r)
	r.Mod(r, n)

	return &SchnorrProof{*r, V}

}

func VerifyProof(proof *SchnorrProof, A ecdsa.PublicKey, UserID []byte) bool {

	//if the key is not on the curve the proof is automatically false
	if !curve.IsOnCurve(A.X, A.Y) {
		return false
	}

	//calculated G * [r] (GxrX represents the X coordinate of G * [r]. Likewise, GxrY represents G * [r]'s Y coordinate)
	GxrX, GxrY := curve.ScalarBaseMult(proof.r.Bytes())

	c := getChallenge(curve.Params(), proof.V.X, A, UserID)

	AxcX, AxcY := curve.ScalarMult(A.X, A.Y, c)

	//finalX represents the addition of points in A * [c] + G * [r]
	finalX, finalY := curve.Add(GxrX, GxrY, AxcX, AxcY)

	return finalX.Cmp(proof.V.X) == 0 && finalY.Cmp(proof.V.Y) == 0

}

func getChallenge(parameters *elliptic.CurveParams, V *big.Int, A ecdsa.PublicKey, UserID []byte) []byte {
	challengeBytes := append(parameters.Gx.Bytes(), V.Bytes()...)
	challengeBytes = append(challengeBytes, A.X.Bytes()...)
	challengeBytes = append(challengeBytes, UserID...)
	challenge := sha256.Sum256(challengeBytes)
	return challenge[:]
}

func SerializeProof(proof *SchnorrProof) []byte {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	encoder.Encode(proof)
	return buf.Bytes()
}

func DeserializeProof(proofBytes []byte) *SchnorrProof {
	var proof SchnorrProof
	decoder := gob.NewDecoder(bytes.NewReader(proofBytes))
	decoder.Decode(&proof)
	return &proof
}
