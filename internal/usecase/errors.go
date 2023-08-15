package usecase

import "errors"

var (
	ErrInternal   = errors.New("internal error")
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)
