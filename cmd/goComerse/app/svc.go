package app

import (
	"database/sql"
	"time"

	authModel "github.com/gocomerse/internal/service/auth"
	"github.com/gocomerse/internal/service/user"
	userModel "github.com/gocomerse/internal/service/user/model"
	userRepo "github.com/gocomerse/internal/service/user/repo"
)

const (
	tokenExpiration = 25 * time.Minute
)

type services struct {
	UserSvc userModel.Service
	AuthSvc *authModel.Service
}

func buildServices(db *sql.DB) *services {
	svc := &services{}
	svc.buildUserService(db)
	svc.buildAuthService()
	return svc
}

func (s *services) buildUserService(db *sql.DB) {
	userRepository := userRepo.NewRepository(db)
	s.UserSvc = user.NewService(userRepository)

}

func (s *services) buildAuthService() {

	s.AuthSvc = authModel.NewService("secret", tokenExpiration, s.UserSvc)
}
