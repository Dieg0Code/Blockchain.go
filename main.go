package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

/*
	A Blockchain is essentially a public database that is distributed accross multiple
	different pairs
*/

//BlockChain : <-
type BlockChain struct {
	blocks []*Block
}

//Block : basic struct
type Block struct {
	Hash     []byte //represent the hash of this block
	Data     []byte //represent the data inside this block
	PrevHash []byte //represent the last block hash, this allow us to link the block together
	// each block inside a blockchain references the last block that was created inside the blockchain
}

//DeriveHash : Generate hash
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// CreateBlock : Create a block
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

// AddBlock : Add a new block to the Blockchain
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

// Genesis : Genesis block
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

//InitBlockChain : <-
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {
	chain := InitBlockChain()

	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("SeconD Block after Genesis")
	chain.AddBlock("Third Block after Genesis")

	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
