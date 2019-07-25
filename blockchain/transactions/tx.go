package transactions

type txOutput struct {
	value         float64
	publicKeyHash []byte
}

type txInput struct {
	ID          []byte
	outputCount int
	sig         []byte
	publicKey   []byte
}
