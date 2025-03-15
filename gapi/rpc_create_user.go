package gapi

import (
	"context"

	backendproduction_2 "github.com/backendproduction-2/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *backendproduction_2.CreateUserRequest) (*backendproduction_2.CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
