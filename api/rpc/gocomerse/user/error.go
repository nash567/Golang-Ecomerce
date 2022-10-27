package user

import "errors"

var (
	ErrInvalidOperation = errors.New("you are not allowed to perform this action")
	ErrIntToStrConv     = errors.New("cannot convert id to int")
	ErrEmptyContext     = errors.New("context is emplty")
)
