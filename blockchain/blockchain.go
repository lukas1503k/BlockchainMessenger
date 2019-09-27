package blockchain
//structure to represent the entire blockchain
import (
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const path string = "/data/ledger"
type Blockchain struct {
	currentLength int
	currentHash   []byte
	db *badger.DB
}

func ledgerExists() bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

func loadBlockchain() *Blockchain {
	if ledgerExists() {
		return loadExistingBlockchain()
	}
	return InitBlockchain()

}
func loadExistingBlockchain() *Blockchain {
	options := badger.Options(nil)
	options.Dir = path
	options.ValueDir = path
	db, err := openDatabase(path, options)
	var newestHash []byte
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("newestHash"))
		handle(err)
		newestHash, err = item.Value()

		return err
	})

	handle(err)
	newBlock := deserializeBlock(newestHash)
	return &Blockchain{newBlock.chainLength, newBlock.blockHash, db}

}

func addNewBlock(chain *Blockchain, messages []*message) {

	//fetch top block using blockchain.currentHash

	newBlock := createBlock(chain.currentHash, messages, chain.currentLength)
	addBlockToChain(newBlock, chain)

}

func getBlockByHash(hash []byte, database *badger.DB) *block {


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
func InitBlockchain() *Blockchain {
	//Precondition, no current blockchain exists

	options := badger.DefaultOptions
	options.Dir = path
	options.ValueDir = path
	var newestHash []byte
	db, err := openDatabase(path, options)
	handle(err)
	err = db.Update(func(txn *badger.Txn) error {
		genesisBlock := createGenesis(nil)
		err = txn.Set(genesisBlock.blockHash, serializeBlock(genesisBlock))
		handle(err)
		err = txn.Set([]byte("newestHash"), genesisBlock.blockHash)

		newestHash = genesisBlock.blockHash

		return err

	})
	handle(err)

	newBlockChain := Blockchain{0, newestHash, db}
	return &newBlockChain

}

func openDatabase(dir string, opts badger.Options) (*badger.DB, error) {
	if db, err := badger.Open(opts); err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := retry(dir, opts); err == nil {
				return db, nil
			}
			log.Println("could not unlock database:", err)
		}
		return nil, err
	} else {
		return db, nil
	}
}

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)
	return db, err
}

func addBlockToChain(newBlock *block, chain *Blockchain) {
	db := chain.db
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("newestBlock"))
		handle(err)
		latestBlock, err := item.Value()
		err = txn.Set(chain.currentHash, latestBlock)
		handle(err)
		err = txn.Set(newBlock.blockHash, serializeBlock(newBlock))
		handle(err)
		chain.currentHash = newBlock.blockHash
		chain.currentLength++
		return err

	})
	handle(err)

}

func handle(err interface{}){
	if err != nil{
		log.Panic(err)
	}
}