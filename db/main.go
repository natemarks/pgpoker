package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/natemarks/secret-hoard/types"
	"github.com/rs/zerolog"
)

// ConnectAndPingDB connects to PostgreSQL and performs a ping to check the connection
func ConnectAndPingDB(secret types.RDSSecretData, log *zerolog.Logger) error {
	// Replace with your actual PostgreSQL connection string
	connectionString := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/postgres?sslmode=disable",
		secret.Username,
		secret.Password,
		secret.Host,
		secret.Port)

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	defer db.Close()

	// Ping the database to check the connection
	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	log.Info().Msgf("Successfully connected to the PostgreSQL database!")
	return nil
}
