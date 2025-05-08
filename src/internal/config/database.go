package config

import (
	"database/sql"

	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

func NewDatabase(
	driverName string,
	connectionURL string,
) *sql.DB {
	db, err := sql.Open(driverName, connectionURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initialize database")
	}

	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")
	}

	return db
}
