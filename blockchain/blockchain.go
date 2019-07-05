package blockchain
//structure to represent the entire blockchain
import(
	badger "github.com/dgraph-io/badger"
	"log"
	"github.com/lukas1503k/msger/blockchain/block"

)
type blockchain struct {
	currentLength int
	currentHash   []byte
	db *badger.DB
}

const dir = "/tmp/ledger"
func addNewBlock(chain *blockchain, messages []*message) {

	//fetch top block using blockchain.currentHash

	newBlock := createBlock(chain.currentHash, messages, chain.currentLength)
	chain.currentHash = newBlock.blockHash
	chain.currentLength = newBlock.chainLength


	blockchain.db.Set()
}


func getBlockByHash(hash []byte ,database badger.DB) *block{


	var searchedBlock *block
	err := database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(hash)
		log.Panic(err)
		serializedBlock, err := item.Value()
		searchedBlock = deserializeBlock(serializedBlock)
		return err
	}
	if err != nil{
		log.Panic(err)
	}
	return searchedBlock

}