package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://neondb_owner:npg_wCMz4j1Epxga@ep-crimson-hill-a12phdeg-pooler.ap-southeast-1.aws.neon.tech/backenddb-2?sslmode=require"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can't connect to the db", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
