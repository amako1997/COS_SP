package storage

import (
	"fmt"
	storageserviceinterface "hidden-service/services/storage/storage-interface"
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
	fmt.Println(storageService.dbPath)
	storageService.options = badger.DefaultOptions
	storageService.options.Dir = storageService.dbPath
	storageService.options.ValueDir = storageService.dbPath
	return storageService
}
func InitializeBadger() badger.Options {
	pwd, _ := os.Getwd()
	storageService := new(Storage)
	storageService.dbPath = pwd + "/database/tmp"
	fmt.Println(storageService.dbPath)
	storageService.options = badger.DefaultOptions
	storageService.options.Dir = storageService.dbPath
	storageService.options.ValueDir = storageService.dbPath
	return storageService.options
}

func (this *Storage) GetDBVolume() (*badger.DB, error) {
	db, err := badger.Open(this.options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger in storage service during volume op")
	}
	return db, nil
}
func (this *Storage) Put(key []byte, val []byte) error {

	db, err := badger.Open(this.options)
	if err != nil {
		return errors.Wrap(err, "failed to open badger in storage service ")
	}
	defer db.Close()

	txn := db.NewTransaction(true)
	defer txn.Discard()

	err = txn.Set(key, val)
	if err != nil {
		return errors.Wrap(err, "failed to save data in badger db ")
	}

	if err := txn.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return nil
}
func (this *Storage) Get(key []byte) ([]byte, error) {

	db, err := badger.Open(this.options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger in storage service ")
	}
	defer db.Close()
	var value []byte
	err = db.View(func(tx *badger.Txn) error {

		item, err := tx.Get(key)
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
