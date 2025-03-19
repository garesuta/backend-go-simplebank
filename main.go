package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/backendproduction-2/api"
	db "github.com/backendproduction-2/db/sqlc"
	_ "github.com/backendproduction-2/doc/statik"
	"github.com/backendproduction-2/gapi"
	backendproduction_2 "github.com/backendproduction-2/pb"
	"github.com/backendproduction-2/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	// run db migration
	runDBMigration(config.MigrationURL, dbSource)

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGRPCServer(config, store)
	// runGinServer(config, store)
}
func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("can't create new migrate instance: ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}

	log.Println("db migrated successfully")
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("can't create server ", err)
	}

	grpcServer := grpc.NewServer()
	log.Println("Registering SimpleBankServer...")
	backendproduction_2.RegisterSimpleBankServer(grpcServer, server)
	log.Println("Reflection enabled...")
	reflection.Register(grpcServer)
	log.Printf("Starting gRPC server on %s...", config.GRPCServerAddress)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("can't create listener: ", err)
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can't start gRPC server ", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("can't create server ", err)
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions:   protojson.MarshalOptions{UseProtoNames: true},
		UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = backendproduction_2.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("can't register handler server ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// fs := http.FileServer(http.Dir("./doc/swagger"))
	// mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("can't create statik file swagger fs ", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("can't create listener: ", err)
	}
	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("can't start HTTP gateway server: ", err)
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
