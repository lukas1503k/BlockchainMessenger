package blockchain

import (
	"crypto/sha256"
	"log"
)

type MerkleTree struct {
	root *MerkleNode
}

type MerkleNode struct {
	data  []byte
	left  *MerkleNode
	right *MerkleNode
}

func getRoot(data [][]byte) []byte {
	root := initTree(data)

	return root.data

}

func initTree(data [][]byte) *MerkleTree {
	// constructs the Merkle Tree for the block
	var nodes []MerkleNode
	//create a node for each piece of data
	for _, datum := range data {
		node := createNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
	if len(nodes) == 0 {
		log.Panic("No nodes created")
	}

	for len(nodes) > 1 {
		// Merkle trees need
		if len(nodes)%2 == 1 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}
		i := 0
		var newLevel []MerkleNode
		for i < len(nodes) {
			newParent := createNode(&nodes[i], &nodes[i+1], nil)
			newLevel = append(newLevel, *newParent)
			i += 2
		}

		nodes = newLevel

	}

	newTree := MerkleTree{&nodes[0]}
	return &newTree

}

func createNode(left, right *MerkleNode, data []byte) *MerkleNode {
	//Helper Function that creates nodes
	if left == nil && right == nil { //if the node is a leaf
		hash := sha256.Sum256(data)
		node := MerkleNode{hash[:], nil, nil}
		return &node
	} else // the node is not a leaf
	{
		//concatenating the two child's data for this node's data
		childData := append(left.data, right.data...)
		hash := sha256.Sum256(childData)
		node := MerkleNode{hash[:], left, right}
		return &node
	}

}
