package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

func UpsertConfig(ctx *context.Context, cfg *ConfigEntity) error {
	// Retrieve the database connection from the context
	db, ok := (*ctx).Value("db").(*sql.DB)
	if !ok || db == nil {
		panic("Database connection not found in context")
	}

	if cfg.Trakt != nil {
		jsonData, err := json.Marshal(cfg.Trakt)
		if err != nil {
			return fmt.Errorf("failed to marshal trakt config: %w", err)
		}
		if err := upsertConfig(ctx, db, "trakt", string(jsonData)); err != nil {
			return err
		}
	}
	if cfg.Emby != nil {
		jsonData, err := json.Marshal(cfg.Emby)
		if err != nil {
			return fmt.Errorf("failed to marshal emby config: %w", err)
		}
		if err := upsertConfig(ctx, db, "emby", string(jsonData)); err != nil {
			return err
		}
	}
	if cfg.Plex != nil {
		jsonData, err := json.Marshal(cfg.Plex)
		if err != nil {
			return fmt.Errorf("failed to marshal plex config: %w", err)
		}
		if err := upsertConfig(ctx, db, "plex", string(jsonData)); err != nil {
			return err
		}
	}
	if cfg.Jellyfin != nil {
		jsonData, err := json.Marshal(cfg.Jellyfin)
		if err != nil {
			return fmt.Errorf("failed to marshal jellyfin config: %w", err)
		}
		if err := upsertConfig(ctx, db, "jellyfin", string(jsonData)); err != nil {
			return err
		}
	}
	if cfg.Cronjob != "" {
		if err := upsertConfig(ctx, db, "cronjob", cfg.Cronjob); err != nil {
			return err
		}

		cronManager, ok := (*ctx).Value("cron").(*CronManager)
		if ok || cronManager != nil {
			// cronManager.UpdateFrequency(ctx, cfg.Cronjob, func() {
			// 	log.Println("Cron job updated with new configuration")
			// })
			// TODO reload cron jobs if cron manager is available
		}
	}

	return nil
}

func upsertConfig(ctx *context.Context, db *sql.DB, cfgType string, cfgData string) error {
	query := `INSERT INTO config (type, data) VALUES (?, ?) 
	ON CONFLICT(type) DO UPDATE SET data = excluded.data`
	_, err := db.ExecContext(*ctx, query, cfgType, cfgData)
	if err != nil {
		log.Printf("Failed to upsert configuration: %v", err)
		return err
	}
	return nil
}

func ReadConfig(ctx *context.Context) (*ConfigEntity, error) {
	// Retrieve the database connection from the context
	db, ok := (*ctx).Value("db").(*sql.DB)
	if !ok || db == nil {
		panic("Database connection not found in context")
	}

	query := `SELECT type, data FROM config`
	rows, err := db.QueryContext(*ctx, query)
	if err != nil {
		log.Printf("Failed to query configurations: %v", err)
		return &ConfigEntity{}, err
	}

	config := ConfigEntity{}

	for rows.Next() {
		var configType, configData string
		if err := rows.Scan(&configType, &configData); err != nil {
			log.Printf("Failed to scan configuration row: %v", err)
			return &ConfigEntity{}, err
		}

		switch configType {
		case "trakt":
			if err := json.Unmarshal([]byte(configData), &config.Trakt); err != nil {
				log.Printf("Failed to unmarshal trakt config: %v", err)
				return &ConfigEntity{}, err
			}
		case "emby":
			if err := json.Unmarshal([]byte(configData), &config.Emby); err != nil {
				log.Printf("Failed to unmarshal emby config: %v", err)
				return &ConfigEntity{}, err
			}
		case "plex":
			if err := json.Unmarshal([]byte(configData), &config.Plex); err != nil {
				log.Printf("Failed to unmarshal plex config: %v", err)
				return &ConfigEntity{}, err
			}
		case "jellyfin":
			if err := json.Unmarshal([]byte(configData), &config.Jellyfin); err != nil {
				log.Printf("Failed to unmarshal jellyfin config: %v", err)
				return &ConfigEntity{}, err
			}
		case "cronjob":
			config.Cronjob = configData
		default:
			log.Printf("Unknown config type: %s", configType)
			return &ConfigEntity{}, fmt.Errorf("unknown config type: %s", configType)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v", err)
		return &ConfigEntity{}, err
	}

	return &config, nil
}
