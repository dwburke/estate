package memory

import (
	"github.com/dwburke/prefs/storage"
)

var MemoryStore map[string]string = map[string]string{}

type Keys struct {
	storage.Storage
}

func New() storage.Storage {
	return &Keys{}
}

func (e *Keys) Get(key string) (string, error) {
	v, ok := MemoryStore[key]
	if !ok {
		return "", storage.ErrNotFound
	}
	return v, nil
}

func (e *Keys) Set(key string, value string) error {
	MemoryStore[key] = value
	return nil
}

func (e *Keys) Delete(key string) error {
	_, ok := MemoryStore[key]
	if ok {
		delete(MemoryStore, key)
	}
	return nil
}

func (e *Keys) Close() {
	return
}
