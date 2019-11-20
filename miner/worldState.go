package miner

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/lukas1503k/msger/blockchain"
	"log"
)

const path string = "/data/ledger"

type WorldState struct {
	db    *badger.DB
	Chain blockchain.Blockchain
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
	db, err := blockchain.OpenDatabase(path, options)
	if err != nil {
		log.Panic(err)
	}

	handle(err)

	state.db = db
	state.chain = *blockchain.InitBlockchain()
	state.exchangesDB = exch

	return state
}

func initAccount(address []byte, balance float64) AccountState {
	return AccountState{
		address:       address,
		accountNounce: 0,
		balance:       balance,
	}

}

func (state *WorldState) addAccountToState(address []byte, balance float64) {
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

func (state *WorldState) VerifyAccountState(address []byte, nounce uint32, amount float64) bool{
	var verification bool
	err := state.db.Update(func (txn *badger.Txn) error{
	accountItem, err := txn.Get(address)
	handle(err)
	var accountBytes []byte
	accountItem.ValueCopy(accountBytes)

	account := DeserializeAccount(accountBytes)
	if account.accountNounce == nounce - 1 && account.balance > amount{
		verification = true
	} else{
		verification = false
		}
	return nil
	})
	return verification

}
func (state *WorldState) UpdateState(address []byte, change float64) {

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


func (state WorldState) GetExchange(exchange blockchain.KeyExchange) (blockchain.KeyExchange, error){

	messageHash := Serialize(exchange)
	var keyExchange blockchain.KeyExchange
	err := state.exchangesDB.Update(func(txn *badger.Txn) error {
		accountItem, err := txn.Get(messageHash)
		var message []byte
		accountItem.ValueCopy(message)
		if err != nil{
			return err
		}
		keyExchange = blockchain.DeserializeKeyExchange(message)
		return nil
	})

	if err != nil{
		return blockchain.KeyExchange{}, err
	}
	return keyExchange, nil

}

func (state WorldState) AddBlock(newBlock blockchain.Block) {
	for i:= 0; i < len(newBlock.messages); i++{
		value := newBlock.messages[i].Amount
		to := newBlock.messages[i].To
		from := newBlock.messages[i].From

		state.UpdateState(to, value)
		state.UpdateState(from, -value)
	}

}
func handle(err interface{}) {
	if err != nil {
		log.Panic(err)
	}

}
