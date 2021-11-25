package post

import "errors"

var (
	ErrNotFound        = errors.New("record not found")
	ErrVersionConflict = errors.New("version conflict")
)
