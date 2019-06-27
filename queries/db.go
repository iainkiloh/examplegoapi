package queries

import "database/sql"

var db *sql.DB

func SetDatabase(database *sql.DB) {
	db = database
}

func Ping() error {
	err := db.Ping()
	return err
}
