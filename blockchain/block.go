package blockchain

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
