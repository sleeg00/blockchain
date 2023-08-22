package main

import (
	"errors"

	"github.com/boltdb/bolt"
)

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next returns next block starting from the tip
func (i *BlockchainIterator) Next() (*Block, error) {

	var block *Block
	var encodedBlock []byte

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock = b.Get(i.currentHash)

		if !(len(encodedBlock) > 0) {

			return errors.New("fail")
		}
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {

		return block, errors.New("fail")
	}

	i.currentHash = block.PrevBlockHash

	return block, nil
}
