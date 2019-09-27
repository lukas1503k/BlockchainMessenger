package miner

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
)

const path string = "/data/ledger"

type WorldState struct {
	db    *badger.DB
	chain blockchain.Blockchain
}

type AccountState struct {
	address       []byte
	accountNounce uint32
	balance       float32
}

func initWorldState() WorldState {
	state := WorldState{}
	options := badger.DefaultOptions(path)
	options.Dir = path
	options.ValueDir = path
	db, err := blockchain.openDatabase(path, options)
	if err != nil {
		log.Panic(err)
	}

	state.db = db
	state.chain = blockchain.InitBlockchain()

	return state
}

func initAccount(address []byte, balance float32) AccountState {
	return AccountState{
		address:       address,
		accountNounce: 0,
		balance:       balance,
	}

}

func (state *WorldState) addAccountToState(address []byte, balance float32) {
	err := state.db.Update(func(txn *badger.Txn) error {
		newAccount := initAccount(address, balance)
		serializedAccount := Serialize(newAccount)
		txn.Set(address, serializedAccount)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func VerifyAccountState(address []byte, nounce int32, amount float32)
func (state *WorldState) updateState(address []byte, change float32) {

	err := state.db.Update(func(txn *badger.Txn) error {
		accountItem, err := txn.Get(address)
		var accountBytes []byte
		accountItem.ValueCopy(accountBytes)
		handle(err)
		account := DeserializeAccount(accountBytes)
		account.balance += change
		txn.Set(address, Serialize(account))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

func Serialize(input interface{}) []byte {
	var output bytes.Buffer

	msgEncoder := gob.NewEncoder(&output)

	err := msgEncoder.Encode(input)
	fmt.Println(err)
	return output.Bytes()
}

func DeserializeAccount(input []byte) AccountState {
	var accountState AccountState

	decoder := gob.NewDecoder(bytes.NewReader(input))

	decoder.Decode(&accountState)

	return accountState
}

func handle(err interface{}) {
	if err != nil {
		log.Panic(err)
	}

}
