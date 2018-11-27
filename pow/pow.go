package pow

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"../tx"
	"../utils"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 8

type IBlock interface {
	GetTimestamp() int64
	GetTransactions() []*tx.Transaction
	GetPrevBlockHash() []byte
	GetHash() []byte
	GetNonce() int
	HashTransactions() []byte
}

type ProofOfWork struct {
	block  IBlock
	target *big.Int
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.GetPrevBlockHash(),
			pow.block.HashTransactions(),
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

	fmt.Println("Mining a new block")

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

func NewProofOfWork(b IBlock) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{block: b, target: target}
	return pow
}
