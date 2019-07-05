package blockchain
//an iterator structure which will allow us to go back in the chain
//key feature needed to the ability to go forward after iterating
//so that nodes can download blocks that they do not have without needing to redownload the whole chain

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/golang-collections/collections/stack"
	"github.com/lukas1503k/msger/blockchain/block"
	"github.com/lukas1503k/msger/blockchain/blockchain"
	"log"

)



type blockIterator struct{
	currentHash []byte
	chain blockchain
	prevHashes stack.Stack
	pos int
	db *badger.DB
}



//Initializes the iterator to be at the newest block
func newIteration(chain blockchain) *blockIterator{

	itr := &blockIterator{chain.currentHash, chain, stack.Stack{}, 0}
	return itr
}



func (iter blockIterator) iterateBackToPos(pos int) *block{
	//get top block from db

	if iter.pos < pos || pos == 0{
		log.Panic("Index Out of Range")
		return nil
	}

	for iter.pos > pos{
		//use helper to iterate



		if pos > iter.chain.currentLength {
			log.Panic("Index Out of Range")
			return nil
		}
	}
	return &block{}//get iter.currenthash block from db
}

func (iter blockIterator)iterForwardToPos(pos int) (*block, error){
	if iter.pos != pos && iter.prevHashes.Len()+iter.pos != pos{
		log.Panic("Index Out of Range")
		return nil, errors.New("Index Out of Range")
	} else{

		for iter.pos != pos{
			iter.getNextBlock()

		}
	}

}

func (iter blockIterator) getNextBlock(){
	iter.currentHash = iter.prevHashes.Pop()
}

func (iter *blockIterator) getPrevBlock(){
	iter.prevHashes.Push(iter.currentHash)
	currentBlock := getBlockByHash(iter.currentHash)
	iter.currentHash = currentBlock.header.prevHash
	currentBlock = nil
	iter.pos += 1

}

func (iter blockIterator) Iterate(pos int) *block{
	if pos == iter.pos{

		return getBlockByHash(iter.currentHash, iter.db)
	} else if pos < iter.pos{
		return iter.iterateBackToPos(pos)
	}
	else{
		block, err := iter.iterForwardToPos(pos)
		errRespond(err)
		return block
	}

}


func errRespond(err error){
	if err != nil{
		fmt.Println(err)
	}
}