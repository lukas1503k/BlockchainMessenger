package blockchain

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/math"
	"log"
	"encoding/gob"
	"math/big"
	"time"

)


const version int = 6
const difficulty int = 1
type block struct {
	header blockHeader
	messages    []*message
	chainLength int
	blockHash []byte
}


type blockHeader struct{
	timestamp time.Time
	MerkleRoot	[]byte
	version int
	Nounce int64
	difficulty int
	prevHash []byte

}

func CreateBlock(prevHash []byte, messages []miner.MessageInterface, chainLength int) *block {



	hashRoot := getHash(messages)

	blockHash, nounce := proofOfWork()


	header := blockHeader{time.Now(), hashRoot, version, nounce, difficulty,
		prevHash}

	newBlock := &block{header, messages, chainLength + 1, blockHash}


	return newBlock

}

func getHash(messages []*message) []byte {

	data := SerializeMessageArray(messages)
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


func (block blockHeader)proofOfWork() ([]byte, int){
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - difficulty))
	nounce := 0
	var hashInt big.Int
	var hash [32]byte
	for nounce <  math.MaxInt64{
		data := combineData(block.merkleRoot, block.prevHash, difficulty, int64(nounce))
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(target) == -1{
			break
		}else{
			nounce++
		}

	}
	return hash[:], nounce

}


func combineData(data []byte, prevHash []byte, difficulty int, nounce int64) []byte{
	combinedData := bytes.Join(
		[][]byte{
			prevHash,
			data,
			IntToBytes(nounce),
			IntToBytes(int64(difficulty)),
		},
		[]byte{})
	return combinedData


}
func IntToBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}

	return buff.Bytes()
}

func (newBlock block) VerifyHash() bool{
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - newBlock.header.difficulty))
	data := combineData(newBlock.merkleRoot, newBlock.prevHash, newBlock.header.difficulty, int64(newBlock.header.nounce))

	var hashInt big.Int

	hashInt.SetBytes(data)
	return hashInt.Cmp(target) == -1
}