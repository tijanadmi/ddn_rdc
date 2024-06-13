package oraclerepo

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/godror/godror"
	"github.com/tijanadmi/ddn_rdc/util"
)

var testDB *sql.DB



func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Panicln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = testDB.PingContext(ctx)
	if err != nil {
		log.Panicln(err)
	}

	testRepo = &OracleDBRepo{DB: testDB}

	// run tests
	code := m.Run()

	os.Exit(code)

}