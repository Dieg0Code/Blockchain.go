package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

/*
	A Blockchain is essentially a public database that is distributed accross multiple
	different pairs
*/

//Block : basic struct
type Block struct {
	Hash         []byte //represent the hash of this block
	Transactions []*Transaction
	PrevHash     []byte //represent the last block hash, this allow us to link the block together
	Nonce        int
	// each block inside a blockchain references the last block that was created inside the blockchain
}

// HashTransactions : provide a unique representation of all of our transactions combined
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions { //iterate through each of the transactions inside of a block
		txHashes = append(txHashes, tx.ID) // append each one to the two dimensional slice of bytes txHashes
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{})) //concatenate all these bytes together and then hash them

	return txHash[:]
}

// CreateBlock : Create a block
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis : Genesis block
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{}) //pase a array of transaction with only the coin base inside of it and an empty solice of bytes
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
