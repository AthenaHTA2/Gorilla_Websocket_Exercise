package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// ConnectDB establishes a connection to the SQLite database
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database/sqlite.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY AUTOINCREMENT, sender TEXT, content TEXT)`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}


