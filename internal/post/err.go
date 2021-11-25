package post

import "errors"

var (
	ErrNotFound        = errors.New("record not found")
	ErrorAlreadyExists = errors.New("record already exists")
)
