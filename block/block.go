package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"

	"../pow"
	"../tx"
)

type Block struct {
	Timestamp     int64
	Transactions  []*tx.Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

/*
func (b *Block) SetHash() {
    timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
    headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
    hash := sha256.Sum256(headers)
    b.Hash = hash[:]
}
*/

func (b *Block) GetTimestamp() int64                { return b.Timestamp }
func (b *Block) GetTransactions() []*tx.Transaction { return b.Transactions }
func (b *Block) GetPrevBlockHash() []byte           { return b.PrevBlockHash }
func (b *Block) GetHash() []byte                    { return b.Hash }
func (b *Block) GetNonce() int                      { return b.Nonce }

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Printf(err.Error())
	}

	return result.Bytes()
}

func newGenesisBlock(coinbase *tx.Transaction) *Block {
	return NewBlock([]*tx.Transaction{coinbase}, []byte{})
}

func NewBlock(txs []*tx.Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), txs, prevBlockHash, []byte{}, 0}
	pow := pow.NewProofOfWork(pow.IBlock(block))
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Printf(err.Error())
	}

	return &block
}
