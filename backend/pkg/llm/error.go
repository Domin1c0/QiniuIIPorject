package llm

import "errors"

var (
	errInvalidRole = errors.New("invalid role in message")
	errNoChoices   = errors.New("no choices returned from llm")
)
