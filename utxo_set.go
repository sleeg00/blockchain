package main

import (
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
)

const utxoBucket = "chainstate"

// UTXOSet represents UTXO set
type UTXOSet struct {
	Blockchain *Blockchain
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions returns the number of transactions in the UTXO set
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex rebuilds the UTXO set
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip of a blockchain
func (u UTXOSet) Update(block Block) {

	log.Println("\n\n블럭이전해시  : ", block.PrevBlockHash)

	db := u.Blockchain.db
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		log.Println("\n\n블럭이전해시1  : ", block.PrevBlockHash)
		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				log.Println("\n\n블럭이전해시 2 : ", block.PrevBlockHash)
				for _, vin := range tx.Vin {
					updatedOuts := TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					if len(outsBytes) == 0 {
						// outsBytes가 비어있는 경우 처리 로직 추가
						continue
					}
					log.Println("\n\n블럭이전해시3  : ", block.PrevBlockHash)
					outs := DeserializeOutputs(outsBytes)
					log.Println("\n\n블럭이전해시4  : ", block.PrevBlockHash)
					for outIdx, out := range outs.Outputs {
						if outIdx != vin.Vout {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}
					log.Println("\n\n블럭이전해시5  : ", block.PrevBlockHash)
					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.Txid)
						if err != nil {
							log.Panic(err)
						}
						log.Println("\n\n블럭이전해시6  : ", block.PrevBlockHash)
					} else {
						err := b.Put(vin.Txid, updatedOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
						log.Println("\n\n블럭이전해시 7 : ", block.PrevBlockHash)
					}

				}
			}
			log.Println("\n\n블럭이전해시8  : ", block.PrevBlockHash)

			newOutputs := TXOutputs{}
			for _, out := range tx.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			log.Println("\n\n블럭이전해시9  : ", block.PrevBlockHash)
			err := b.Put(tx.ID, newOutputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
			log.Println("\n\n블럭이전해시 10 : ", block.PrevBlockHash)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	log.Println("\n\n블럭이전해시 11 : ", block.PrevBlockHash)
}

func (u UTXOSet) UpdateTx(txo *Transaction) {

	db := u.Blockchain.db
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		if !txo.IsCoinbase() {

			for _, vin := range txo.Vin {
				updatedOuts := TXOutputs{}
				outsBytes := b.Get(vin.Txid)
				if len(outsBytes) == 0 {
					// outsBytes가 비어있는 경우 처리 로직 추가
					continue
				}
				outs := DeserializeOutputs(outsBytes)

				for outIdx, out := range outs.Outputs {
					if outIdx != vin.Vout {
						updatedOuts.Outputs = append(updatedOuts.Outputs, out)
					}
				}

				if len(updatedOuts.Outputs) == 0 {
					err := b.Delete(vin.Txid)
					if err != nil {
						log.Panic(err)
					}
				} else {
					err := b.Put(vin.Txid, updatedOuts.Serialize())
					if err != nil {
						log.Panic(err)
					}
				}

			}

			newOutputs := TXOutputs{}
			for _, out := range txo.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			err := b.Put(txo.ID, newOutputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
