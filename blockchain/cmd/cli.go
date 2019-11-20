package cmd

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/lukas1503k/msger/blockchain"
	"log"
	"os"
	"runtime"
	"strconv"
)
func createBlockChain(){
	chain := blockchain.InitBlockchain()
	fmt.Print("Chain Created")
}
var account wallet.Account

func initMessageExchange(to string){
	toAddress, err := hex.DecodeString(to)
	if err != nil{
		log.Panic(err)
	}
	wallet.InitExchange(account, toAddress)
	sendMessage
}

func resptoMessageExchange(from string){
	break
}


func initAccount(){



}