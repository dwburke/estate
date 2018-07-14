package storage

import (
	"errors"
)

var (
	ErrNotFound = errors.New("storage: record not found")
)

func NotFound(e error) bool {
	if e == ErrNotFound {
		return true
	}
	return false
}

type Storage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Close()
}
