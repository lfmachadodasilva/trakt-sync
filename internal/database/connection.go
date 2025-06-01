package database

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

func Connect() *sql.DB {
	once.Do(func() {
		var err error
		// Ensure the folder exists
		if _, err := os.Stat("./data"); os.IsNotExist(err) {
			if err := os.Mkdir("./data", os.ModePerm); err != nil {
				log.Fatalf("Failed to create data folder: %v", err)
			}
		}
		// Ensure the file exists
		if _, err := os.Stat("./data/data.db"); os.IsNotExist(err) {
			file, err := os.Create("./data/data.db")
			if err != nil {
				log.Fatalf("Failed to create database file: %v", err)
			}
			file.Close()
		}
		dbInstance, err = sql.Open("sqlite3", "./data/data.db")
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Create the table if it does not exist
		createTableQuery := `CREATE TABLE IF NOT EXISTS config (
			type TEXT PRIMARY KEY,
			data TEXT
		)`
		_, err = dbInstance.Exec(createTableQuery)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	})
	return dbInstance
}
