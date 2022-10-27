package auth

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	logModel "github.com/gocomerse/internal/logger/model"
)

type Interceptor struct {
	enabled bool
	log     logModel.Logger
	authSvc Service
}

func NewAuthInterceptor(log logModel.Logger, enabled bool, service Service) *Interceptor {
	return &Interceptor{
		log:     log,
		authSvc: service,
		enabled: enabled,
	}
}

func (i *Interceptor) Auth(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	if info.FullMethod == "/user.UserService/Login" {
		return handler(ctx, req)
	}
	// Extract token from header context
	idToken, err := ExtractToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to fetch token")
	}
	claims, err := i.authSvc.Verify(strings.TrimPrefix(idToken[0], "Bearer "))
	if err != nil {
		return nil, err
	}

	if info.FullMethod == "/user.UserService/UpdateUser" {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md.Append("userId", strconv.Itoa(claims.ID))

		}
		newCtx := metadata.NewIncomingContext(ctx, md)

		return handler(newCtx, req)
	}
	// var expired bool

	// if expired {
	// 	if i.auth.Verify()(ctx, i.log, strings.TrimPrefix(idToken[0], "Bearer ")) {
	// 		return handler(ctx, req)
	// 	}
	// }

	return handler(ctx, req)
}

func ExtractToken(ctx context.Context) ([]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	idToken := md.Get("authorization")
	if len(idToken) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "token not submitted")
	}
	return idToken, nil
}
