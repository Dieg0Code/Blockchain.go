package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte //hash
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	ID  []byte //references the transaction that the output is inside of
	Out int    //index of the output
	Sig string // is a script witch provides the data witch is used in the outputs Pubkey
	// but in this impl the Sig is going to be the users account
	//inputs are just references to previous outputs
}

type TxOutput struct {
	Value  int    //value in tokens
	Pubkey string //needed to unlock tokens inside value field

	//Outputs are indivisible you can't reference a part of an output
	/*
		An analogy to this is :
		if you walk into a store and you buy 5 dolars then pay with 10 dollars to the cashier
		the cashier cant just rip the 10 dollar bill on half and hand you the other half
		you back, instead he have to give you back a 5 dollar bill.

		So there is 10 tokens inside of our output we need to create new outputs one with
		5 tokens inside of it and another one with 5 tokens inside of it.

		Esp:

		- Los Outputs son indivisibles, no puedes hacer referencia a solo una parte de este.

		Una analogia :
		Si entras a una tienda y compras 5 dolares en productos luego le pagas al cajero con un
		billete de 10 dolares el cajero no puede solo partir el billete por la mitad y luego darte
		la otra mitad, en vez de eso tiene que darte de vuelta otro billete de 5 dolares

		Asique si hay 10 tokens dentro de nuestro output necesitariamos crear nuevos outputs, uno con 5
		tokens dentro y otro con 5 tokens dentro también.
	*/
}

/*
	In our blockchain we have our genesis block in that block we
	also have our first transaction this is wath's called a coinbase transaction,
	in this transaction we only have one input and only one output and the input inside of it
	rether than referencing an older output, because there's no olders outputs, just references an
	empty output it also doesn't store a signature instead store a bunch of arbitrary data, the coinbase
	also has what's called the subsidy or a reward attached to it, this reward is released to a single account
	when that individual mines the coinbase.
	In our impl we are just going to add a constant for our coinbases then we're mainly just doing this
	to make things simpler for now.

	Esp:

	En nuestra blockchain tenemos el bloque génesis en ese bloque también tenemos la primera transacción,
	esto es lo que es llamado transacción "Coinbase", en esta transacción solo tenemos un solo input y un solo output
	y dentro de ese output en vez de hacer referencia a un output anterior, porque no hay outputs anteriores, solo
	hace referencia a un output vació, tampoco almacena una firma en vez de eso almacena un puñado de datos arbitrarios,
	la coinbase también tiene lo que es llamado el subsidio o una recompensa añadida a el, esta recompensa es liberada
	a una cuenta individual cuando esa cuenta mina la coinbase.
	En nuestra implementación solo vamos a añadir una constante para nuestra coinbase luego manualmente vamos a hacer eso
	para hacer la cosas mas simples por ahora.
*/

//SetID : Creates a hash based on bytes that represents the transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx) // encode the transaction
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes()) //hash the bytes portion of our encoded bytes
	tx.ID = hash[:]                       // set that hash into the transaction id
}

// CoinbaseTx :
func CoinbaseTx(to, data string) *Transaction {
	if data == "" { //empty
		data = fmt.Sprint("Coins to %s", to)
	}
	txin := TxInput{[]byte{}, -1, data} //empty slice of bytes for id, outIndex = -1, signature
	txout := TxOutput{100, to}          //reward, pubkey string for this output as a reference to the "to" address

	//Instance of the transaction struct
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}} //nil for id, inp, out
	tx.SetID()                                                 //create hash id for this transaction

	return &tx //return a reference for this transaction
}

// NewTransaction : creates a new transaction
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount { // check if the amount is greater than the accumulator
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs { // iterate through  valid outputs
		txID, err := hex.DecodeString(txid) // take each and decode the string txid into bytes
		Handle(err)

		for _, out := range outs { // iterate through the outs
			input := TxInput{txID, out, from} // creat a new input for each of the unspent outputs
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to}) // create an output with the amount that we are going to send and then the "to" address which is the person that we are sending to

	if acc > amount { // check if the amount is less than the accumulated which means that the amount that the from user has is greater than the amount that he's trying to send
		outputs = append(outputs, TxOutput{acc - amount, from}) // create a second output. Is created if there is any left over tokens in the original sender account
	}

	tx := Transaction{nil, inputs, outputs} // instancies a transaction and passed an inputs and an outputs
	tx.SetID()                              //set the id of the transaction

	return &tx // return the reference to the transaction
}

//IsCoinbase : Allow us tho determine if a transaction is a coninbase transaction or not
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1 // if all is true is a coinbase transaction
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data // check if is equals to the signature attached to the input struct
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.Pubkey == data //check if is equal to the pubkey if is true means that the account(data) owns information inside the output
}
