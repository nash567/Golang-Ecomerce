package auth

import "errors"

var (
	ErrBadTokenSignMethod = errors.New("unexpected token signing method")
	ErrInvalidClaims      = errors.New("invalid token claims")
	ErrMoreRecordFound    = errors.New("more than one records found with same info")
)
