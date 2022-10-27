package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	logModel "github.com/gocomerse/internal/logger/model"
	userPB "github.com/gocomerse/internal/pb/gocomerse/user"
	"github.com/gocomerse/internal/service/auth"
	"github.com/gocomerse/internal/service/user/model"
)

type Server struct {
	userPB.UnimplementedUserServiceServer
	log     logModel.Logger
	userSvc model.Service
	authSvc auth.Service
}

func NewServer(userSvc model.Service, log logModel.Logger, authSvc auth.Service) *Server {
	return &Server{log: log, userSvc: userSvc, authSvc: authSvc}

}
func (s *Server) GetUser(ctx context.Context, req *userPB.GetUserRequest) (*userPB.GetUserResponse, error) {

	var queryParams model.QueryParams
	getQueryParams(&queryParams, req)
	users := make([]*userPB.User, 0)
	res, err := s.userSvc.Get(ctx, s.log.WithFields(logModel.Fields{
		"server": "User",
		"method": "GetUser",
		"rid":    uuid.New(),
	}), queryParams, false)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot list users:%s", err)
	}
	for _, user := range res {
		newUser := &userPB.User{
			Id:        int64(user.UserID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
		}

		users = append(users, newUser)
	}
	return &userPB.GetUserResponse{
		Users: users,
	}, nil
}

func (s *Server) GetUserByID(ctx context.Context, req *userPB.UserId) (*userPB.User, error) {

	user, err := s.userSvc.GetByID(ctx, s.log.WithFields(logModel.Fields{
		"server": "User",
		"method": "GetUserByID",
		"rid":    uuid.New(),
	}), int(req.GetId()))

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot find user with  id: %d", req.GetId())
	}

	return &userPB.User{
		Id:        int64(user.UserID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}, nil
}

func (s *Server) RegisterUser(ctx context.Context, req *userPB.RegisterUserRequest) (*userPB.RegisterUserResponse, error) {

	user, err := s.userSvc.Create(ctx, s.log.WithFields(logModel.Fields{
		"server": "User",
		"method": "Register User",
		"rid":    uuid.New(),
	}), model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
	})

	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot create user: %v", err)
	}

	return &userPB.RegisterUserResponse{
		Id:        int64(user.UserID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *userPB.UpdateUserResquest) (*userPB.UpdateUserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%w", ErrEmptyContext)
	}

	updateID, err := strconv.Atoi(md["userId"][0])
	if err != nil {
		return nil, fmt.Errorf("%w", ErrIntToStrConv)
	}
	if updateID != int(req.Id) {
		return nil, fmt.Errorf("%w", ErrInvalidOperation)
	}
	user, err := s.userSvc.Update(ctx, s.log.WithFields(logModel.Fields{
		"server": "User",
		"method": "Update User",
		"rid":    uuid.New(),
	}), model.User{
		UserID:    int(req.GetId()),
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Phone:     req.GetPhone(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot update user")
	}

	return &userPB.UpdateUserResponse{
		UserId:    int64(user.UserID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *userPB.UserId) (*userPB.DeleteUserResponse, error) {
	err := s.userSvc.Delete(ctx, s.log.WithFields(logModel.Fields{
		"server": "User",
		"method": "Delete User",
		"rid":    uuid.New(),
	}), int(req.GetId()))

	if err != nil {
		return &userPB.DeleteUserResponse{
			Status:  "Failure",
			Message: "cannot delete user",
		}, status.Errorf(codes.Aborted, "cannot delete the user")
	}
	return &userPB.DeleteUserResponse{
		Status:  "success",
		Message: "user deleted successfully",
	}, nil
}

func (s *Server) Login(ctx context.Context, req *userPB.LoginRequest) (*userPB.LoginResponse, error) {
	user, err := s.authSvc.Login(ctx, s.log, model.UserCredential{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &userPB.LoginResponse{
		User: &userPB.User{
			Id:        int64(user.UserID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
		},
		Token: user.Token,
	}, nil
}

func getQueryParams(query *model.QueryParams, req *userPB.GetUserRequest) {
	if req.FirstName != nil {

		query.FirstName = *req.FirstName
	}
	if req.LastName != nil {

		query.LastName = *req.LastName
	}
	if req.Email != nil {

		query.Email = *req.Email
	}

	if req.Order != nil {

		query.Order = *req.Order
	}
	if req.Sort != nil {

		query.Sort = *req.Sort
	}

	if req.Limit != nil {

		query.Limit = int(*req.Limit)
	}
	if req.Page != nil {

		query.Page = int(*req.Page)
	}

}
