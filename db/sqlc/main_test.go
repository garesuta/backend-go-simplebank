package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/backendproduction-2/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can't load config", err)
	}

	dbDriver := config.DBDriver
	dbSource := config.DBSource

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can't connect to the db", err)
	}
	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}
