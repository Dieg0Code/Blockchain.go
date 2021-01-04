package blockchain

import (
	"bytes"
	"crypto/sha256"
)

/*
	A Blockchain is essentially a public database that is distributed accross multiple
	different pairs
*/

//BlockChain : <-
type BlockChain struct {
	Blocks []*Block
}

//Block : basic struct
type Block struct {
	Hash     []byte //represent the hash of this block
	Data     []byte //represent the data inside this block
	PrevHash []byte //represent the last block hash, this allow us to link the block together
	// each block inside a blockchain references the last block that was created inside the blockchain
}

//DeriveHash : Generate the hash of a Block
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
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

// Genesis : Genesis block
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

//InitBlockChain : <-
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
