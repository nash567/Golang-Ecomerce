package model

import (
	"context"

	logModel "github.com/gocomerse/internal/logger/model"
)

type user interface {
	Get(ctx context.Context, log logModel.Logger, queryParams QueryParams, pass bool) ([]*User, error)
	GetByID(ctx context.Context, log logModel.Logger, id int) (*User, error)
	Create(ctx context.Context, log logModel.Logger, user User) (*User, error)
	Update(ctx context.Context, log logModel.Logger, user User) (*User, error)
	Delete(ctx context.Context, log logModel.Logger, id int) error
}
type Service interface {
	user
}

type Repository interface {
	user
}
