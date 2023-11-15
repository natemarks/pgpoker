package db

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/natemarks/secret-hoard/types"
	"github.com/rs/zerolog"
)

func TCPConn(data types.RDSSecretData) (err error) {
	conn, err := net.DialTimeout("tcp", data.Host+":"+strconv.Itoa(data.Port), 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

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
