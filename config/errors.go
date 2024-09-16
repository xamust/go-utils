package config

import "errors"

var (
	ErrNotPtr       = errors.New("dest must be pointer")
	ErrPathNotFound = errors.New("path is empty")
)
