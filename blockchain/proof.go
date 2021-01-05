package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

/*

Consensu algorithms or proof algorithms:

Proof of work:
Secure the blockchain by forcing the network to do work to add a block
to the chain. This idea of work is essentially just computational power, when
you hear miners "mining bitcoin" they are doing this proof of work algorithm so they can
sign the blocks on the blockchain and the reason why they get fees is because they're
essentially powering the network by running the proof of work algorithm and also by doing that they
make the actual block and the data inside this block is more secure. Proof of work also goes hand and hand
with validation of this proof there is the say when a user does the work to sign a block they the need to provide
proof of this work. An important concept they goes hand and hand with this is the fact that the work must be hard
to do but that proving that work must be relatively easy.
*/

// Take the data from the block

// Create a counter (nonce) which starts at 0
//Increments upwards theoretically infinitely

// Create a hash of the data plus the counter

// Check the hash to see if it meets a set of requirements
/* this is where the idea of difficulty comes in, if the hash meets
the set of requirements then we use that hash and we say that it sign the block
otherwise we go back and we create another hash an we repeat this proccess until
we get a hash that does meet the set of requirements
*/

// Requirements:
// The first few bytes must contains 0s
/* In the original BitCoin proof of work specification wish is called
"HashCash" the original difficulty was to get 20 consecutive bits of the hash as 0s
this requirement get adjusted over time, and that is essentially the difficulty so the
difficulty goes up and that means that we must have more proceding 0s in front of that hash
for it to be valid
*/

// Difficulty : stay static
const Difficulty = 12

/* In our implementation of this algorithm our difficulty is going to stay
static.
In a real block chain generally you would have an algorithm that with slowly increment this difficulty
over a large period of time, the main reason you want to do this is to account for the increasing number
of miners of the network and also account for the increase in computation power of computers in general, because
we wanna make the time to mine a block stay the same and we also want to have the block rate stay the same, this
means that you need to have a certain amount of computational power running the proof of work algorithm to
produce blocks at that rate but also keep the time time to sign a block down
*/

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

/* what makes this proof of work algorithm very secure for a blockchain is this idea that if you
want to change a block inside of the chain you have to recalculate the hash it self wish takes a
large amount of time and the you have to recalculate every block after that block hash as well to be
able to say "this is valid data" so even now we can validate that a block is valid fairly quickely
the actual work needed to create to create the block or to sign the block is pretty difficult and that
protects the blockchain data from being tempered by one single entity or even by a large amounts of entities
*/

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
