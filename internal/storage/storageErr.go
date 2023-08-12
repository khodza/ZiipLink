package storage

import "errors"

var (
	ErrURLNotFound   = errors.New("url not found")
	ErrUrlExists     = errors.New("url exists")
	ErrAliasExists   = errors.New("alias exists")
	ErrAliasNotFound = errors.New("alias not found")
)
