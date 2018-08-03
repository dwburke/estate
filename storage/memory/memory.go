package memory

import (
	"errors"

	"github.com/dwburke/prefs/storage/common"
	"github.com/dwburke/prefs/storage/meta"
)

var ErrNotFound = errors.New("storage: record not found")

var MemoryStore map[string][]byte = map[string][]byte{}

type Storage struct {
	meta.Storage
}

func New() (*Storage, error) {
	return &Storage{}, nil
}

func (e *Storage) Get(key string) ([]byte, error) {
	v, ok := MemoryStore[key]
	if !ok {
		return nil, common.ErrNotFound
	}
	return v, nil
}

func (e *Storage) Set(key string, value []byte) error {
	MemoryStore[key] = value
	return nil
}

func (e *Storage) Delete(key string) error {
	_, ok := MemoryStore[key]
	if ok {
		delete(MemoryStore, key)
	}
	return nil
}

func (e *Storage) Close() error {
	return nil
}
