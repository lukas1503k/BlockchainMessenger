package blockchain

import (
	"bytes"
	"encoding/gob"
)

type block struct {
	prevHash    []byte
	hash        []byte
	activity    []*message
	chainLength int
}

func createBlock(prevHash []byte, messages []*message, chainLength int) *block {

	hashRoot := getHash(prevHash, messages)
	newBlock := &block{prevHash, hashRoot, messages, chainLength + 1}

	return newBlock

}

func getHash(prevHash []byte, messages []*message) []byte {

	data := append(prevHash, messages.Serialize()...)

	return merkletree.getRoot(data)
}

func createGenesis() *block {
	genesis := &block{0, nil, nil, 0}

	return genesis
}

func serializeBlock(block *block) []byte {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	encoder.Encode(block)

	return buf.Bytes()

}

func deserializeBlock(incBlock []byte) *block {

	var decodedBlock block
	decoder := gob.NewDecoder(bytes.NewReader(incBlock))
	decoder.Decode(&decodedBlock)

	return &decodedBlock
}
