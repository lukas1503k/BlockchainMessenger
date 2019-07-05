package blockchain

import (
	"bytes"
	"log"
	"encoding/gob"
	"time"

)


const version int = 6
const difficulty int = 0
type block struct {
	header blockHeader
	messages    []*message
	chainLength int
	blockHash []byte
}


type blockHeader struct{
	timestamp time.Time
	merkleRoot	[]byte
	version int
	nounce int
	difficulty int
	prevHash []byte

}

func createBlock(prevHash []byte, messages []*message, chainLength int) *block {



	hashRoot := getHash(messages)

	blockHash, nounce := proofOfWork()


	header := blockHeader{time.Now(), hashRoot, version, nounce, difficulty,
		prevHash}

	newBlock := &block{header, messages, chainLength + 1, blockHash}

	return newBlock

}

func getHash(messages []*message) []byte {

	data := serializeMessageArray(messages)

	return getRoot(data)
}

func createGenesis(messages []*message) *block {

	genesis := createBlock(nil, messages, 0)

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


func proofOfWork() ([]byte, int){

	return *new([]byte), 0
}