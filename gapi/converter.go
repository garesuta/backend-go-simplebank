package gapi

import (
	db "github.com/backendproduction-2/db/sqlc"
	backendproduction_2 "github.com/backendproduction-2/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *backendproduction_2.User {
	return &backendproduction_2.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
