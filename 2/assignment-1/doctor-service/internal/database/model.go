package database

import "database/sql"

type Model interface {
	TableName() string
	TableCreateQuery() string
}

func SqlInitModels(db *sql.DB, models []Model) error {
	for _, model := range models {
		query := model.TableCreateQuery()

		_, err := db.Exec(query)

		if err != nil {
			return err
		}
	}

	return nil
}
