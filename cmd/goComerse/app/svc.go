package app

import (
	"context"
	"database/sql"

	logModel "github.com/gocomerse/internal/logger/model"
	"github.com/gocomerse/internal/service/user"
	userModel "github.com/gocomerse/internal/service/user/model"
	userRepo "github.com/gocomerse/internal/service/user/repo"
)

type services struct {
	UserSvc userModel.Service
}

func buildServices(ctx context.Context, log logModel.Logger, db *sql.DB) *services {
	svc := &services{}
	svc.buildUserService(db)
	return svc
}

func (s *services) buildUserService(db *sql.DB) {
	userRepository := userRepo.NewRepository(db)
	s.UserSvc = user.NewService(userRepository)

}
