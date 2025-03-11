package main

import (
	"database/sql"
	"log"

	"github.com/backendproduction-2/api"
	db "github.com/backendproduction-2/db/sqlc"
	"github.com/backendproduction-2/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config ", err)
	}

	dbDriver := config.DBDriver
	dbSource := config.DBSource
	serverAddress := config.ServerAddress

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can't connect to the db ", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can't create server ", err)
	}

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("can't start the server ", err)
	}
}
