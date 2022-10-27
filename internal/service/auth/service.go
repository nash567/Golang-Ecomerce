package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	logModel "github.com/gocomerse/internal/logger/model"
	"github.com/gocomerse/internal/service/user/model"
)

type Service struct {
	// log           logModel.Logger
	secretKey     string
	tokenDuration time.Duration
	userSvc       model.Service
}
type Claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	ID    int    `json:"id"`
}

func NewService(secretKey string, tokenDuration time.Duration, userSvc model.Service) *Service {
	return &Service{secretKey: secretKey, tokenDuration: tokenDuration, userSvc: userSvc}
}

func (s *Service) GenerateToken(user *model.UserCredential) (*string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenDuration).Unix(),
		},
		Email: user.Email,
		ID:    user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, fmt.Errorf("error creating signed token: %w", err)
	}
	return &signedToken, nil
}

func (s *Service) Verify(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("%w", ErrBadTokenSignMethod)
			}

			return []byte(s.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("%w", ErrInvalidClaims)
	}

	return claims, nil
}

func (s *Service) Login(ctx context.Context, log logModel.Logger, user model.UserCredential) (*model.User, error) {
	users, err := s.userSvc.Get(ctx, log, model.QueryParams{Email: user.Email}, true)
	if err != nil {
		log.WithError(err).Errorf("user not found by this email")
		return nil, fmt.Errorf("user not found by this email: %w", err)
	}
	if len(users) > 1 {
		log.WithError(err).Errorf("more than one user found by this email")

		return nil, fmt.Errorf("%w", ErrMoreRecordFound)

	}
	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(user.Password))
	if err != nil {
		log.WithError(err).Errorf("unauthorized user! enter correct password: ")

		return nil, fmt.Errorf("unauthorized user! enter correct password: %w", err)
	}

	token, err := s.GenerateToken(&model.UserCredential{
		Email:    user.Email,
		Password: user.Password,
		ID:       users[0].UserID,
	})
	if err != nil {
		log.WithError(err).Errorf("error generating token:")

		return nil, fmt.Errorf("error generating token: %w", err)
	}
	users[0].Token = *token
	return users[0], nil
}
