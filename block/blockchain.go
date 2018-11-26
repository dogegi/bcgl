package block

import (
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "bcgl.db"
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	DB  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	newBlock := NewBlock(data, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b != nil {
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			err = b.Put([]byte("l"), newBlock.Hash)
			bc.tip = newBlock.Hash

			if err != nil {
				log.Printf(err.Error())
			}
		}

		return nil
	})

	if err != nil {
		log.Printf(err.Error())
	}
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.DB}
	return bci
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := newGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)

			if err != nil {
				log.Printf(err.Error())
			}

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Printf(err.Error())
	}

	bc := Blockchain{tip, db}
	return &bc
}

func newGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
