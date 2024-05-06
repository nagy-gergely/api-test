package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqlite() *sql.DB {
	db, err := sql.Open("sqlite3", "./test.db")

	if err != nil {
		log.Fatal(err)
	}

	initDb(db)

	return db
}

func initDb(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
}
