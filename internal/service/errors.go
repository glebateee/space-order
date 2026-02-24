package service

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrInvalid  = errors.New("invalid request")
)
