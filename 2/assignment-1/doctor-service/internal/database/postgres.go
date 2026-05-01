package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func PostgresConnectDB(connection string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connection)

	if err != nil {
		return nil, err
	}

	return db, nil
}
