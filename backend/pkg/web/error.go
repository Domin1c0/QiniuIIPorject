package web

import "errors"

var (
	ErrInternalError        = errors.New("internal server error")
	ErrRequestMissingFields = errors.New("request missing fields")
	ErrRequestInvalid       = errors.New("request invalid")
	ErrNotFound             = errors.New("resource not found")
	ErrPermissionDenied     = errors.New("forbidden")
	ErrUserOrPasswordWrong  = errors.New("invalid username or wrong password")
)
