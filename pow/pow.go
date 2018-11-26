package pow

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"../utils"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 24

type IBlock interface {
	GetTimestamp() int64
	GetData() []byte
	GetPrevBlockHash() []byte
	GetHash() []byte
	GetNonce() int
}

type ProofOfWork struct {
	block  IBlock
	target *big.Int
}

func NewProofOfWork(b IBlock) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{block: b, target: target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.GetPrevBlockHash(),
			pow.block.GetData(),
			utils.IntToHex(pow.block.GetTimestamp()),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.GetData())

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		//fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.GetNonce())
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
