package usecase

import "errors"

// list of general errors on usecase layer
var (
	ErrInternal   = errors.New("internal error")
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)
