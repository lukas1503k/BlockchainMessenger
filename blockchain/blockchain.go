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
const exchanges string = "/data/exchanges"

type Blockchain struct {
	currentLength int
	currentHash   []byte
	db            *badger.DB
	exchangesDB   *badger.DB
}

func ledgerExists() bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}

	return true
}

func LoadBlockchain() *Blockchain {
	if ledgerExists() {
		return loadExistingBlockchain()
	}
	return InitBlockchain()

}
func loadExistingBlockchain() *Blockchain {
	options := badger.Options(nil)
	options.Dir = path
	options.ValueDir = path
	db, err := OpenDatabase(path, options)
	var newestHash []byte
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("newestHash"))
		handle(err)
		item.ValueCopy(newestHash)

		return err
	})

	exch, err := OpenDatabase(exchanges, options)

	handle(err)
	newBlock := deserializeBlock(newestHash)
	return &Blockchain{newBlock.chainLength, newBlock.blockHash, db, exch}

}

func AddNewBlock(chain *Blockchain, messages []*Message) {

	//fetch top block using blockchain.currentHash

	newBlock := CreateBlock(chain.currentHash, messages, chain.currentLength)
	AddBlockToChain(newBlock, chain)

}

func getBlockByHash(hash []byte, database *badger.DB) *block {

	var searchedBlock *block
	err := database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(hash)
		log.Print(err)

		err = item.Value(func(val []byte) error {
			searchedBlock = deserializeBlock(val)
			return nil
		})

		return err
	})
	if err != nil {
		log.Panic(err)
	}
	return searchedBlock

}
func InitBlockchain() *Blockchain {
	//Precondition, no current blockchain exists

	options := badger.DefaultOptions(path)
	options.Dir = path
	options.ValueDir = path
	var newestHash []byte
	db, err := OpenDatabase(path, options)
	handle(err)
	err = db.Update(func(txn *badger.Txn) error {
		genesisBlock := createGenesis(nil)
		err = txn.Set(genesisBlock.blockHash, serializeBlock(genesisBlock))
		handle(err)
		err = txn.Set([]byte("newestHash"), genesisBlock.blockHash)

		newestHash = genesisBlock.blockHash

		return err

	})
	exch, err := OpenDatabase(exchanges, options)

	handle(err)

	newBlockChain := Blockchain{1, newestHash, db, exch}
	return &newBlockChain

}

func OpenDatabase(dir string, opts badger.Options) (*badger.DB, error) {
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

func AddBlockToChain(newBlock *block, chain *Blockchain) {
	db := chain.db
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("newestBlock"))
		handle(err)
		var latestBlock []byte
		item.ValueCopy(latestBlock)
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

func handle(err interface{}) {
	if err != nil {
		log.Panic(err)
	}
}
