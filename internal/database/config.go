package database

import (
	"database/sql"
	"log"
)

func UpsertConfig(db *sql.DB, configType string, configData string) error {
	query := `INSERT INTO config (type, data) VALUES (?, ?) 
	ON CONFLICT(type) DO UPDATE SET data = excluded.data`
	_, err := db.Exec(query, configType, configData)
	if err != nil {
		log.Printf("Failed to upsert configuration: %v", err)
		return err
	}
	return nil
}

func ReadConfig(db *sql.DB, configType string) (string, error) {
	query := `SELECT data FROM config WHERE type = ?`
	var configData string
	err := db.QueryRow(query, configType).Scan(&configData)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Printf("Failed to read configuration: %v", err)
		return "", err
	}
	return configData, nil
}
