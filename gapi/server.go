package gapi

import (
	"fmt"

	db "github.com/backendproduction-2/db/sqlc"
	backendproduction_2 "github.com/backendproduction-2/pb"
	"github.com/backendproduction-2/token"
	"github.com/backendproduction-2/util"
)

// Server serves gRPC requests for our banking service
type Server struct {
	backendproduction_2.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new HTTP server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker}

	return server, nil
}
