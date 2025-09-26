package web

import "errors"

var (
	ErrRequestMissingFields = errors.New("request missing fields")
	ErrUserOrPasswordWrong  = errors.New("invalid username or wrong password")
)
