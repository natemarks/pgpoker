package db

import (
	"database/sql"
	//required to embed the sql query
	_ "embed"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	// required for sql access
	_ "github.com/lib/pq"
	"github.com/natemarks/secret-hoard/types"
)

//go:embed embed/list_roles.sql
var listRolesQuery string

// TCPConn tests a TCP connection to the RDS instance
func TCPConn(data types.RDSSecretData) (err error) {
	conn, err := net.DialTimeout("tcp", data.Host+":"+strconv.Itoa(data.Port), 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// SecretToConn given an RDS secret for an RDS instance and a database name, return a connection
// if the dbname is ”, the returned session is connected to an instance, but not a database.
// otherwise, the session will try to connect to the specific database
func SecretToConn(secret types.RDSSecretData, dbname string) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=disable",
		secret.Host,
		strconv.Itoa(secret.Port),
		secret.Username,
		secret.Password)
	if dbname != "" {
		dataSourceName = dataSourceName + fmt.Sprintf(" dbname=%s", dbname)
	}

	db, err := sql.Open(
		"postgres", dataSourceName)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

// ListRoles runs the SQL query to list roles in the PostgreSQL database
func ListRoles(conn *sql.DB) (roles []string, err error) {

	rows, err := conn.Query(listRolesQuery)
	if err != nil {
		return roles, err
	}

	for rows.Next() {
		var roleName string
		err := rows.Scan(&roleName)
		if err != nil {
			return roles, err
		}
		roles = append(roles, roleName)
	}

	return roles, nil
}

// CheckInstance Check the instance configuration
func CheckInstance(secret types.RDSSecretData, log *zerolog.Logger) {
	instanceLog := log.With().Str("host", secret.Host).Logger()
	conn, err := SecretToConn(secret, "postgres")
	if err != nil {
		instanceLog.Error().Err(err).Msgf("error opening DB connection to master secret host: %v", secret.Host)
		return
	}
	instanceLog.Info().Msgf("Opened DB connection to master secret host: %v", secret.Host)

	defer conn.Close()
	roles, err := ListRoles(conn)
	if err != nil {
		instanceLog.Error().Err(err).Msgf("error listing roles on master secret host: %v", secret.Host)
	}
	instanceLog.Info().Msgf("roles on %v: %v", secret.Host, roles)

}
