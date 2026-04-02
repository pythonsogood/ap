package database

type Model interface {
	TableName() string
	TableCreateQuery() string
}
