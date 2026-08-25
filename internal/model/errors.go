package model

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid request")
	ErrAcked    = errors.New("delivery already acknowledged")
)
