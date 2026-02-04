package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/godror/godror"
	"github.com/tijanadmi/ddn_rdc/cmd/api"
	db "github.com/tijanadmi/ddn_rdc/repository"
	"github.com/tijanadmi/ddn_rdc/repository/oraclerepo"
	"github.com/tijanadmi/ddn_rdc/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := openDB(config.DBDriver, config.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	dbStore := &oraclerepo.OracleDBRepo{DB: conn}
	runGinServer(config, dbStore)
}

func runGinServer(config util.Config, store db.DatabaseRepo) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}

func openDB(dbDriver string, dns string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dns)
	if err != nil {
		return nil, err
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
