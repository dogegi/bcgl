package block

import (
	"log"

	"github.com/boltdb/bolt"
)

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get([]byte(i.currentHash))
		block = DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		log.Printf(err.Error())
	}

	i.currentHash = block.PrevBlockHash
	return block
}
