package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

//BadgerDB v1.5.4
// Badger is a key - value database written in pure Go
const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

// BlockChain : BlockChain struct
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

//DBexists : allow us to determinate if the badgerDB exists
func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

//InitBlockChain : initialize the DB and the blockchain as well
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func ContinueBlockChain(address string) *BlockChain {
	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one!!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()

		return err
	})

	Handle(err)

	chain := BlockChain{lastHash, db}

	return &chain
}

//AddBlock : add a new block to the chain
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()

		return err
	})
	Handle(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.Value()
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}

//FindUnspentTransactions : Find all unspent transactions witch are assigned to an address
func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction //array of transactions

	spentTXOs := make(map[string][]int) // where keys are strings and values are slice of ints

	iter := chain.Iterator() // iterate through the blockchain in the data base

	for {
		block := iter.Next()

		for _, tx := range block.Transactions { // iterate through each of the transactions inside of a block
			txID := hex.EncodeToString(tx.ID) // take the transaction IDs of each of the transactions and encoded into hexadecimal

		Outputs: // label to break from the inside to this portion of the for loop and not break up of these other two for loops
			for outIdx, out := range tx.Outputs { //iterate through of the outputs inside of the transaction
				if spentTXOs[txID] != nil { //check if the output is inside of our map
					for _, spentOut := range spentTXOs[txID] { //if is not iterate through the map
						if spentOut == outIdx { // if spendOut is equal to the output index
							continue Outputs // continues the Outputs for loop
						}
					}
				}
				if out.CanBeUnlocked(address) { // determinate if the output can be unlocked by the address that we are searching for
					unspentTxs = append(unspentTxs, *tx) //take each of the transactions that can be unlock by these address an put in into the unspent transactions
				}
			}
			if tx.IsCoinbase() == false { // check if the transaction is a coinbase transaction or not
				for _, in := range tx.Inputs { // if not iterate through the transaction inputs, is a way to find other outputs that are referenced by inputs
					if in.CanUnlock(address) { // if we find other outputs check if we can unlock that outputs with the address
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out) // if we can put it insede of the map
					}
				}
			}
		}

		if len(block.PrevHash) == 0 { // if the block is the genesis block break
			break
		}
	}

	return unspentTxs //return the unspent transactions array which would have all unspent transaction which are assigned to the user account which we push through this function
}

// FindUTXO : find all the unspent transactions outputs
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions { // iterate through the unspent transactions
		for _, out := range tx.Outputs { // iterate through the outputs in the transactions
			if out.CanBeUnlocked(address) { // check if the outputs can be unlocked by the address
				UTXOs = append(UTXOs, out) // if they can we add them to the unspent transactions output var
			}
		}
	}

	return UTXOs // the return the unspent transactions
}

// FindSpendableOutputs : find all the unspent outputs and then ensure they have enough tokens inside of them
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)                // unspent outputs
	unspentTxs := chain.FindUnspentTransactions(address) // unspent transactions
	accumulated := 0

Work:
	for _, tx := range unspentTxs { // iterate through the unspent transactions
		txID := hex.EncodeToString(tx.ID) // Encode transaction id into hexadecimal and assign it to txID

		for outIdx, out := range tx.Outputs { // iterate through the outputs inside of the unspent transactions
			if out.CanBeUnlocked(address) && accumulated < amount { // check if the output can be unlocked by the address and if the accumulated is less than the amount that we want to send
				accumulated += out.Value                              // increment the accumulated value by the output value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx) // add the output index to the output map

				if accumulated >= amount {
					break Work
				}

			}

		}
	}

	return accumulated, unspentOuts
}
