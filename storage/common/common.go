package common

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("storage: record not found")
	ErrInvalidDatabase = errors.New("storage: invalid database type")
)
