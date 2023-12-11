package db

import (
	"database/sql"
	"log"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "Planner.db")
	if err != nil {
		log.Fatal(err)
	}

	db.Exec(`
		CREATE TABLE if not exists planner (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Title VARCHAR(255),
		Description VARCHAR(255),
		Status INTEGER,
		Priority VARCHAR(255)
		);
	`)

	return db
}
