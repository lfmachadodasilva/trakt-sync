package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"trakt-sync/internal/models"
)

func UpsertConfig(config *models.Config) error {
	db := GetAndConnect()

	if config.Trakt != nil {
		jsonData, err := json.Marshal(config.Trakt)
		if err != nil {
			return fmt.Errorf("failed to marshal trakt config: %w", err)
		}
		if err := upsertConfig(db, "trakt", string(jsonData)); err != nil {
			return err
		}
	}
	if config.Emby != nil {
		jsonData, err := json.Marshal(config.Emby)
		if err != nil {
			return fmt.Errorf("failed to marshal emby config: %w", err)
		}
		if err := upsertConfig(db, "emby", string(jsonData)); err != nil {
			return err
		}
	}
	if config.Plex != nil {
		jsonData, err := json.Marshal(config.Plex)
		if err != nil {
			return fmt.Errorf("failed to marshal plex config: %w", err)
		}
		if err := upsertConfig(db, "plex", string(jsonData)); err != nil {
			return err
		}
	}
	if config.Jellyfin != nil {
		jsonData, err := json.Marshal(config.Jellyfin)
		if err != nil {
			return fmt.Errorf("failed to marshal jellyfin config: %w", err)
		}
		if err := upsertConfig(db, "jellyfin", string(jsonData)); err != nil {
			return err
		}
	}

	return nil
}

func upsertConfig(db *sql.DB, configType string, configData string) error {
	query := `INSERT INTO config (type, data) VALUES (?, ?) 
	ON CONFLICT(type) DO UPDATE SET data = excluded.data`
	_, err := db.Exec(query, configType, configData)
	if err != nil {
		log.Printf("Failed to upsert configuration: %v", err)
		return err
	}
	return nil
}

func ReadConfig() (models.Config, error) {
	db := GetAndConnect()

	query := `SELECT type, data FROM config`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to query configurations: %v", err)
		return models.Config{}, err
	}

	config := models.Config{}

	for rows.Next() {
		var configType, configData string
		if err := rows.Scan(&configType, &configData); err != nil {
			log.Printf("Failed to scan configuration row: %v", err)
			return models.Config{}, err
		}

		switch configType {
		case "trakt":
			if err := json.Unmarshal([]byte(configData), &config.Trakt); err != nil {
				log.Printf("Failed to unmarshal trakt config: %v", err)
				return models.Config{}, err
			}
		case "emby":
			if err := json.Unmarshal([]byte(configData), &config.Emby); err != nil {
				log.Printf("Failed to unmarshal emby config: %v", err)
				return models.Config{}, err
			}
		case "plex":
			if err := json.Unmarshal([]byte(configData), &config.Plex); err != nil {
				log.Printf("Failed to unmarshal plex config: %v", err)
				return models.Config{}, err
			}
		case "jellyfin":
			if err := json.Unmarshal([]byte(configData), &config.Jellyfin); err != nil {
				log.Printf("Failed to unmarshal jellyfin config: %v", err)
				return models.Config{}, err
			}
		default:
			log.Printf("Unknown config type: %s", configType)
			return models.Config{}, fmt.Errorf("unknown config type: %s", configType)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v", err)
		return models.Config{}, err
	}

	return config, nil
}
