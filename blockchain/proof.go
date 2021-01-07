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

Consensus algorithms or proof algorithms:

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

EspVersion:

Proof of work:

Este algoritmo asegura la blockchain forzando a la red a hacer "trabajo" para añadir un bloque nuevo a la cadena.
Esta idea de "trabajo" es en esencia solo poder computacional, cuando oyes hablar sobre mineros "minando bitcoin"
lo que estos mineros están haciendo es correr el algoritmo proof of work asi ellos pueden firmar los bloques de
la blockchain, la razón por la que los mineros reciben recompensas es porque ellos básicamente están dando
poder de computo a la red al correr el algoritmo proof of work, ademas al hacer esto hacen que el bloque y la data
que contiene el bloque sea mas segura. Proof of work también va de la mano con la validación de la "prueba", es decir
cuando un usuario hace el trabajo para firmar el bloque necesita proveer pruebas de su trabajo. Un concepto
importante que va de la mano con esto es el hecho de que el trabajo debe ser difícil de hacer, pero probar que se hizo
el trabajo deberia ser relativamente fácil.
*/

// Take the data from the block

// Create a counter (nonce) which starts at 0
//Increments upwards theoretically infinitely

// Create a hash of the data plus the counter

// Check the hash to see if it meets a set of requirements
/* this is where the idea of difficulty comes in, if the hash meets
the set of requirements then we use that hash and we say that it sign the block
otherwise we go back and we create another hash and we repeat this proccess until
we get a hash that does meet the set of requirements


Esp:

- checa el hash para ver si cumple un conjunto de requisitos.
Aquí es donde la idea de dificultad sale a la luz, si el hash cumple los requisitos
entonces usamos ese hash y decimos que firme el bloque, de otra manera volvemos y creamos otro
hash, repetimos este proceso hasta que obtengamos un hash que cumpla con los requisitos
*/

// Requirements:
// The first few bytes must contains 0s
/* In the original BitCoin proof of work specification wish is called
"HashCash" the original difficulty was to get 20 consecutive bits of the hash as 0s
this requirement get adjusted over time, and that is essentially the difficulty so the
difficulty goes up and that means that we must have more proceding 0s in front of that hash
for it to be valid

Esp:

- Requisitos:
Los primeros bytes deben contener 0s
En la especificación original de proof of work de Bitcoin la cual es llamada "HashCash"
la dificultad original era tener 20 bits consecutivos como 0s, este requisito iba siendo
ajustado sobre el tiempo, eso es esencialmente la dificultad, que la dificultad vaya en incremento
significa que progresivamente deben haber mas 0s en el principio del hash para que este sea valido
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
produce blocks at that rate but also keep the time to sign a block down

Esp:
En nuestra implementación de este algoritmo la dificultad se va a mantener estática.
En una blockchain real generalmente deberías de tener un algoritmo el cual vaya lentamente
incrementando la dificultad durante un largo periodo de tiempo, la razon principal de esto es que tienes que
tener en cuenta el numero creciente de mineros de la red y ademas tener en cuenta el incremento de poder
computacional en general, porque queremos hacer que el tiempo de minado de un bloque permanezca igual ademas de
la taza de creación de estos, esto significa que necesitas tener cierta cantidad de poder de computo alimentando el
algoritmo de proof of work para producir bloques a la misma taza pero ademas mantener igual el tiempo que toma
firmar un bloque

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
large amount of time and then you have to recalculate every block after that block hash as well to be
able to say "this is valid data" so even now we can validate that  block is valid fairly quickely
the actual work needed to create the block or to sign the block is pretty difficult and that
protects the blockchain data from being tempered by one single entity or even by a large amounts of entities

Esp:
Lo que hace este algoritmo muy seguro para una blockchain es la idea de que si quieres cambiar
un bloque dentro de la cadena tiendrias que recalcular el hash mismo lo cual toma un largo periodo
de tiempo y luego tendrias que recalcular cada hash de cada bloque concecutivo a ese para poder decir
"esta data es valida" asique incluso ahora que podemos validar un bloque relativamente facil el trabajo
necesario para crear un bloque y firmarlo es bastante dificil, eso es lo que protege a la informacion
contenida en la blockchain de ser moderada por una sola entidad o por un conjunto de entidades
*/

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
