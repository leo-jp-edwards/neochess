package store

import "errors"

var (
	ErrNotFound  = errors.New("store: not found")
	ErrDuplicate = errors.New("store: duplicate")
	ErrConflict  = errors.New("store: conflict")
)
