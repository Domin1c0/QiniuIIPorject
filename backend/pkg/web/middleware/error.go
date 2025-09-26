package middleware

import (
	"errors"
	"net/http"
)

var (
	ErrMissingToken = errors.New("missing token")
	ErrInvalidToken = errors.New("invalid token")
)

func Error(w http.ResponseWriter, status int, err error) {
	http.Error(w, err.Error(), status)
}
