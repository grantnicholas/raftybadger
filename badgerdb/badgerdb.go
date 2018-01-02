package badgerdb

import (
	"github.com/dgraph-io/badger"
)

type BadgerDB struct {
	*badger.DB
}

func GetDB() *BadgerDB {
	opts := badger.DefaultOptions
	opts.SyncWrites = true
	opts.Dir = "/tmp/badger"
	opts.ValueDir = "/tmp/badger"
	db, err := badger.Open(opts)
	if err != nil {
		panic("Error while opening badger database")
	}
	badgerTypedDB := BadgerDB{db}
	return &badgerTypedDB
}

func (db *BadgerDB) GetValue(key string) (string, error) {
	value := ""

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		value = string(val)
		return nil
	})

	return value, err
}

func (db *BadgerDB) SetValue(key, value string) (err error) {
	db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
	return err
}
