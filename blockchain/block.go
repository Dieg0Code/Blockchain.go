package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

/*
	A Blockchain is essentially a public database that is distributed accross multiple
	different pairs
*/

//Block : basic struct
type Block struct {
	Hash     []byte //represent the hash of this block
	Data     []byte //represent the data inside this block
	PrevHash []byte //represent the last block hash, this allow us to link the block together
	Nonce    int
	// each block inside a blockchain references the last block that was created inside the blockchain
}

// CreateBlock : Create a block
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis : Genesis block
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

//BadgerDb Serialize - Deserialize

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
