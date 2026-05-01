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
