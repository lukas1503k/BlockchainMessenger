package blockchain

type blockchain struct {
	currentLength int
	currentHash   []byte
}

func addNewBlock(chain *blockchain, messages []*message) {

	//fetch top block using blockchain.currentHash

	newBlock := createBlock(chain.currentHash, messages, chain.currentLength)

	chain.currentHash = newBlock.hash
	chain.currentLength = newBlock.chainLength

	//add new block to badger db

}
