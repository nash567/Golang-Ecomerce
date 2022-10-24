package user

import (
	// logModel "github.com/gocomerse/internal/logger/model".

	userPB "github.com/gocomerse/internal/pb/gocomerse/user"
)

type Server struct {
	userPB.UnimplementedUserServiceServer
	// logger logModel.Logger
}

// func NewServer(server models.Server) *Server {;}
