package db

import (
	"github.com/boltdb/bolt"
	"github.com/guiwoo/gocoin/utils"
)

const (
	dbName      = "blockchain.db"
	dataBucket  = "data"
	blockBucket = "blocks"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		//initialized db
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucket([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucket([]byte(blockBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		blocksBucket := t.Bucket([]byte(blockBucket))
		err := blocksBucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockChain(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		dataBucket := t.Bucket([]byte(dataBucket))
		err := dataBucket.Put([]byte("checkpoint"), data)
		return err
	})
	utils.HandleErr(err)
}
