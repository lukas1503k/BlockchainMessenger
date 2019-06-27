package blockchain

import(
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type wallet struct{
	publicKey []byte
	privateKey ecdsa.PrivateKey
}

func createWallet() *wallet{
	//creates a wallet
	privKey, pubKey = getKeys()

	newWallet := wallet{pubKey, privKey}

	return &newWallet
}

func getKeys() (ecdsa.PrivateKey, []byte) {
	// creates a unique key pair for the wallet
	curve := elliptic.P256()
	newPrivateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	newPublicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *newPrivateKey, newPublicKey
	

}

fun