package database

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
)

func SQLiteConnectDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func SQLiteInitDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS doctors (
		id TEXT UNIQUE NOT NULL,
		full_name TEXT NOT NULL,
		specialization TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	)
	`

	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func NewSQLiteDB(dataSourceName string) (*sql.DB, error) {
	db, err := SQLiteConnectDB(dataSourceName)

	if err != nil {
		return nil, err
	}

	if err := SQLiteInitDB(db); err != nil {
		return nil, err
	}

	return db, nil
}
