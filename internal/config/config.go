package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"trakt-sync/internal/database"
)

func UpsertConfig(config *ConfigEntity) error {
	db := database.GetAndConnect()

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

func ReadConfig() (ConfigEntity, error) {
	db := database.GetAndConnect()

	query := `SELECT type, data FROM config`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to query configurations: %v", err)
		return ConfigEntity{}, err
	}

	config := ConfigEntity{}

	for rows.Next() {
		var configType, configData string
		if err := rows.Scan(&configType, &configData); err != nil {
			log.Printf("Failed to scan configuration row: %v", err)
			return ConfigEntity{}, err
		}

		switch configType {
		case "trakt":
			if err := json.Unmarshal([]byte(configData), &config.Trakt); err != nil {
				log.Printf("Failed to unmarshal trakt config: %v", err)
				return ConfigEntity{}, err
			}
		case "emby":
			if err := json.Unmarshal([]byte(configData), &config.Emby); err != nil {
				log.Printf("Failed to unmarshal emby config: %v", err)
				return ConfigEntity{}, err
			}
		case "plex":
			if err := json.Unmarshal([]byte(configData), &config.Plex); err != nil {
				log.Printf("Failed to unmarshal plex config: %v", err)
				return ConfigEntity{}, err
			}
		case "jellyfin":
			if err := json.Unmarshal([]byte(configData), &config.Jellyfin); err != nil {
				log.Printf("Failed to unmarshal jellyfin config: %v", err)
				return ConfigEntity{}, err
			}
		default:
			log.Printf("Unknown config type: %s", configType)
			return ConfigEntity{}, fmt.Errorf("unknown config type: %s", configType)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v", err)
		return ConfigEntity{}, err
	}

	return config, nil
}
