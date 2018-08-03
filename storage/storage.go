package storage

import (
	"github.com/spf13/viper"

	"github.com/dwburke/lode/storage/common"
	"github.com/dwburke/lode/storage/memory"
	"github.com/dwburke/lode/storage/meta"
	"github.com/dwburke/lode/storage/mysql"
)

type Storage struct {
	engine meta.Storage
}

func New() (*Storage, error) {

	storage_type := viper.GetString("lode.storage.type")

	var engine meta.Storage
	var err error

	switch storage_type {
	case "memory":
		engine, err = memory.New()
	case "mysql":
		engine, err = mysql.New()
	default:
		err = common.ErrInvalidDatabase
	}

	if err != nil {
		return nil, err
	}

	return &Storage{engine}, nil
}

func (store *Storage) Set(key string, value []byte) error {
	return store.engine.Set(key, value)
}

func (store *Storage) Get(key string) ([]byte, error) {
	return store.engine.Get(key)
}

func (store *Storage) Delete(key string) error {
	return store.engine.Delete(key)
}

func (store *Storage) Close() error {
	return store.engine.Close()
}
