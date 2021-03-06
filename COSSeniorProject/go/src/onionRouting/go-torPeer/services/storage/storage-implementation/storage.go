package storage

import (
	"fmt"
	storageserviceinterface "onionRouting/go-torPeer/services/storage/storage-interface"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
)

type Storage struct {
	dbPath  string
	options badger.Options
}

func NewStorage() storageserviceinterface.StorageService {
	pwd, _ := os.Getwd()
	storageService := new(Storage)
	storageService.dbPath = pwd + "/database/tmp"
	storageService.options = badger.DefaultOptions
	storageService.options.Dir = storageService.dbPath
	storageService.options.ValueDir = storageService.dbPath
	return storageService
}

func (this *Storage) Put(key string, val []byte) error {

	db, err := badger.Open(this.options)
	if err != nil {
		return errors.Wrap(err, "failed to open badger in storage service ")
	}
	defer db.Close()
	err = db.Update(func(tx *badger.Txn) error {

		err := tx.Set([]byte(key), val)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to save data in badger db ")
	}
	return nil
}
func (this *Storage) Get(key string) ([]byte, error) {

	db, err := badger.Open(this.options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger in storage service ")
	}
	defer db.Close()
	var value []byte
	err = db.View(func(tx *badger.Txn) error {

		item, err := tx.Get([]byte(key))
		if err != nil {
			return err
		}
		itemBytes, err := item.ValueCopy(nil)
		if err != nil {
			return errors.Wrap(err, "failed to read the value from badger ")
		}
		fmt.Println("data", string(itemBytes))
		value = itemBytes
		return nil
	})

	if err != nil {
		return nil, err
	}
	return value, nil
}
