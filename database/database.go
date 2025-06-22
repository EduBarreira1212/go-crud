package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Conect() (*sql.DB, error) {
	dsn := "user:user123@tcp(localhost:3306)/go-crud-db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
