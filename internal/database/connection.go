package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbInstance     *sql.DB
	once           sync.Once
	databasePath   = "./data/database.db"
	dataFolderPath = "./data"
)

func GetAndConnect(ctx *context.Context) *sql.DB {
	once.Do(func() {
		var err error
		// Ensure the folder exists
		if _, err := os.Stat(dataFolderPath); os.IsNotExist(err) {
			if err := os.Mkdir(dataFolderPath, os.ModePerm); err != nil {
				log.Fatalf("Failed to create data folder: %v", err)
			}
		}
		// Ensure the file exists
		if _, err := os.Stat(databasePath); os.IsNotExist(err) {
			file, err := os.Create(databasePath)
			if err != nil {
				log.Fatalf("Failed to create database file: %v", err)
			}
			file.Close()
		}
		dbInstance, err = sql.Open("sqlite3", databasePath)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	})
	return dbInstance
}
