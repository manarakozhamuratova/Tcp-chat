package repository

import (
	"database/sql"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "forum.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	err = Migrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrations(db *sql.DB) error {
	openTables, err := os.Open("./migrations/init.sql")
	if err != nil {
		return err
	}
	allTables, err := io.ReadAll(openTables)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(allTables))
	if err != nil {
		return err
	}
	return nil
}
