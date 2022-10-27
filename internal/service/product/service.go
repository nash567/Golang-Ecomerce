package product

import (
	"context"
	"fmt"

	logModel "github.com/gocomerse/internal/logger/model"
	"github.com/gocomerse/internal/service/product/model"
)

type Service struct {
	repo model.Repository
}

func NewService(repo model.Repository) *Service {

	return &Service{repo: repo}

}

func (s *Service) Get(ctx context.Context, logger logModel.Logger) ([]*model.Product, error) {
	res, err := s.repo.Get(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("error fetching user %w", err)
	}
	return res, nil
}

func (s *Service) Create(ctx context.Context, logger logModel.Logger, product model.Product) (*model.Product, error) {

	res, err := s.repo.Create(ctx, logger, product)

	if err != nil {
		return nil, fmt.Errorf("error Creating user %w", err)
	}
	return res, nil
}

func (s *Service) Update(ctx context.Context, logger logModel.Logger, product model.Product) (*model.Product, error) {
	res, err := s.repo.Update(ctx, logger, product)

	if err != nil {
		return nil, fmt.Errorf("error updatimg product %w", err)
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, logger logModel.Logger, id int) error {
	err := s.repo.Delete(ctx, logger, id)

	if err != nil {
		return fmt.Errorf("error updatimg product %w", err)
	}
	return nil
}
