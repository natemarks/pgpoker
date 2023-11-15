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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		secret.Host,
		secret.Port,
		secret.Username,
		secret.Password,
		"postgres")
	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", psqlInfo)
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
