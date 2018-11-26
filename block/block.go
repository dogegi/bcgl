package block

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"../pow"
)

type Block struct {
	Timestamp     int64
	Data          []byte
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

func (b *Block) GetTimestamp() int64      { return b.Timestamp }
func (b *Block) GetData() []byte          { return b.Data }
func (b *Block) GetPrevBlockHash() []byte { return b.PrevBlockHash }
func (b *Block) GetHash() []byte          { return b.Hash }
func (b *Block) GetNonce() int            { return b.Nonce }

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Printf(err.Error())
	}

	return result.Bytes()
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
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
