package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cgf mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cgf.FormatDSN())
	if err != nil {
		log.Fatal(err)

	}

	return db, nil
}
