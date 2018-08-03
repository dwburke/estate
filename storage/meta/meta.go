package meta

import (
//"github.com/spf13/viper"
)

//var (
//ErrNotFound = errors.New("storage: record not found")
//)

//func NotFound(e error) bool {
//if e == ErrNotFound {
//return true
//}
//return false
//}

type Storage interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Close() error
}
