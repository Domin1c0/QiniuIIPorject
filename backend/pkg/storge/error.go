package storage

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidUserID = errors.New("invalid user id")
)
