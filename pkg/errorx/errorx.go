package errorx

import "errors"

var (
	ErrOrderValidation = errors.New("error of validation order")
	ErrOrderNotFound   = errors.New("order not found")
	ErrInternal        = errors.New("internal error")
)
