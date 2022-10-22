package db

import (
	"database/sql"
	"fmt"
)

const (
	maxAllowedPacket = 104857600 // 100MB
	txIsolation      = "READ-COMMITTED"
)

// NewConnection returns a database connection for the specified database. If the
// specified database is not one of the know database names, or if any error occurs
// while creating the database connection, this function panics. Connections are
// created only once (the first time this function is called with a given parameter
// value); subsequent calls return the same connection.

func NewConnection(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Name)

	connDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open db connection: %w", err)
	}

	err = connDB.Ping()

	if err != nil {
		return nil, fmt.Errorf("unable to ping to db: %w", err)
	}

	return connDB, nil
}

// Verify ensures the connection is available for use.
func Verify(conn *sql.DB) error {
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}
