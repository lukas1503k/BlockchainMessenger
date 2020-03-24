package main

import (
	"encoding/hex"
	"fmt"
	"github.com/lukas1503k/BlockchainMessenger/wallet"
	"github.com/lukas1503k/msger/blockchain"
	"log"
)

var chain *blockchain.Blockchain

func createBlockChain() {
	chain = blockchain.LoadBlockChain()
	fmt.Print("Chain Created")
}

var account wallet.Account

func initMessageExchange(to string) {
	toAddress, err := hex.DecodeString(to)
	if err != nil {
		log.Panic(err)
	}
	wallet.InitExchange(account, toAddress)
	//sendMessage
}

func resptoMessageExchange(from string) {
	return
}

func initAccount() {
	return

}

func main() {
	chain := blockchain.InitBlockChain()

	blockchain.AddNewBlock(chain, nil)
	fmt.Printf("Previous Hash: %x\n", chain.currentHash)

	/*
		createChain := flag.String("Start", "", "Creates the blockchain or loads one")

		if(createChain){
			createBlockChain()
		}
	*/

}
