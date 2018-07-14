package memory

import (
	"github.com/dwburke/prefs/storage"
)

type Keys struct {
	store map[string]string
}

func New() storage.Storage {
	return &Keys{store: map[string]string{}}
}

func (e *Keys) Get(key string) (string, error) {
	v, ok := e.store[key]
	if !ok {
		return "", storage.ErrNotFound
	}
	return v, nil
}

func (e *Keys) Set(key string, value string) error {
	e.store[key] = value
	return nil
}

func (e *Keys) Delete(key string) error {
	return nil
}

func (e *Keys) Close() {
	return
}
