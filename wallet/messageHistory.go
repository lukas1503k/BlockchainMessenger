package wallet

type MessageKeyPair struct {
	message []blockchain.Message
	key     []byte
}

type MessageChain struct {
	contactID    []byte
	messages     []MessageKeyPair
	messageCount int
	//ephemeralKey ecdsa.PrivateKey
}

func AddMessageToChain(newMessage message)
