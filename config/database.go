package config

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

// LoadDatabase initializes a database connection
func LoadDatabase(ctx context.Context) (*sql.DB, error) {
	dbConfig := GetDBConfig()

	db, err := sql.Open("mysql", dbConfig.ConnectionString())
	if err != nil {
		log.Error().Err(err).Msg("Failed to open connection to database")
		return nil, err
	}

	// Optional: Test the database connection with context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Optional: Test the database connection
	if err := db.PingContext(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	log.Info().Msg("Connection to database successful")
	return db, nil
}
