package user

import (
	"context"
	"fmt"

	logModel "github.com/gocomerse/internal/logger/model"
	userHelper "github.com/gocomerse/internal/service/user/helper"
	"github.com/gocomerse/internal/service/user/model"
)

type Service struct {
	repo model.Repository
}

func NewService(repo model.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, logger logModel.Logger, queryParams model.QueryParams, pass bool) ([]*model.User, error) {
	res, err := s.repo.Get(ctx, logger, queryParams, pass)
	if err != nil {
		return nil, fmt.Errorf("error fetching user %w", err)
	}
	return res, nil
}

func (s *Service) GetByID(ctx context.Context, logger logModel.Logger, id int) (*model.User, error) {
	res, err := s.repo.GetByID(ctx, logger, id)

	if err != nil {
		return nil, fmt.Errorf("error fetching user %w", err)
	}
	return res, nil
}

func (s *Service) Create(ctx context.Context, logger logModel.Logger, user model.User) (*model.User, error) {

	err := userHelper.GenerateHash(&user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash :%w", err)
	}
	res, err := s.repo.Create(ctx, logger, user)

	if err != nil {
		return nil, fmt.Errorf("error Creating user %w", err)
	}
	return res, nil
}

func (s *Service) Update(ctx context.Context, logger logModel.Logger, user model.User) (*model.User, error) {
	res, err := s.repo.Update(ctx, logger, user)

	if err != nil {
		return nil, fmt.Errorf("error updatimg user %w", err)
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, logger logModel.Logger, id int) error {
	err := s.repo.Delete(ctx, logger, id)

	if err != nil {
		return fmt.Errorf("error updatimg user %w", err)
	}
	return nil
}

func (s *Service) Login(ctx context.Context, log logModel.Logger, user model.UserCredential) error {
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
	return nil
}
