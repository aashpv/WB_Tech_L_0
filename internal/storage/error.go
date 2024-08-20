package storage

import "errors"

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrConnect        = errors.New("connect error")
	ErrInternalServer = errors.New("internal server error")
)
