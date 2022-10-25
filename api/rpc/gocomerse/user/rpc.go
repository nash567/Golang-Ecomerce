package user

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logModel "github.com/gocomerse/internal/logger/model"
	userPB "github.com/gocomerse/internal/pb/gocomerse/user"
	"github.com/gocomerse/internal/service/user/model"
)

type Server struct {
	userPB.UnimplementedUserServiceServer
	log logModel.Logger
	svc model.Service
}

func NewServer(svc model.Service, log logModel.Logger) *Server {
	return &Server{log: log, svc: svc}

}
func (s *Server) GetUser(ctx context.Context, req *userPB.GetUserRequest) (*userPB.GetUserResponse, error) {
	users := make([]*userPB.User, 0)
	res, err := s.svc.Get(ctx, s.log.WithFields(logModel.Fields{
		"server": "User",
		"method": "GetUser",
		"rid":    uuid.New(),
	}))
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
	user, err := s.svc.GetByID(ctx, s.log.WithFields(logModel.Fields{
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

func (s *Server) RegisterUser(ctx context.Context, req *userPB.RegisterUserRequest) (*userPB.User, error) {
	user, err := s.svc.Add(ctx, s.log.WithFields(logModel.Fields{
		"server": "User",
		"method": "Register User",
		"rid":    uuid.New(),
	}), model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	})

	if err != nil {
		return nil, status.Errorf(codes.Aborted, "cannot create user: %v", err)
	}

	return &userPB.User{
		Id:        int64(user.UserID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *userPB.UpdateUserResquest) (*userPB.UpdateUserResponse, error) {
	user, err := s.svc.Update(ctx, s.log.WithFields(logModel.Fields{
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
	err := s.svc.Delete(ctx, s.log.WithFields(logModel.Fields{
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
