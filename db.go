package main

import (
	"database/sql"
	"log"
)

// InitDB create the database table and indexes
func InitDB(dsn string) *sql.DB {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal("error db")
	}
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS annotations (
		id INTEGER PRIMARY KEY, 
		annoid VARCHAR, 
		created_at DATETIME, 
		target VARCHAR, 
		manifest VARCHAR, 
		body TEXT)`)
	statement.Exec()

	statement, _ = db.Prepare("CREATE UNIQUE INDEX IF NOT EXISTS annotation_id ON annotations (annoid);")
	statement.Exec()
	return db
}
