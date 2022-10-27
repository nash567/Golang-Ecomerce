package model

import (
	"context"

	logModel "github.com/gocomerse/internal/logger/model"
)

type product interface {
	Get(ctx context.Context, log logModel.Logger) ([]*Product, error)
	Create(ctx context.Context, log logModel.Logger, user Product) (*Product, error)
	Update(ctx context.Context, log logModel.Logger, user Product) (*Product, error)
	Delete(ctx context.Context, log logModel.Logger, id int) error
}
type Service interface {
	product
}

type Repository interface {
	product
}
