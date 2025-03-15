package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/backendproduction-2/api"
	db "github.com/backendproduction-2/db/sqlc"
	"github.com/backendproduction-2/gapi"
	backendproduction_2 "github.com/backendproduction-2/pb"
	"github.com/backendproduction-2/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config ", err)
	}

	dbDriver := config.DBDriver
	dbSource := config.DBSource
	// serverAddress := config.ServerAddress

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can't connect to the db ", err)
	}
	store := db.NewStore(conn)
	runGRPCServer(config, store)
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("can't create server ", err)
	}

	grpcServer := grpc.NewServer()
	backendproduction_2.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("can't create listener ", err)
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can't start gRPC server ", err)
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can't create server ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("can't start the server ", err)
	}
}
