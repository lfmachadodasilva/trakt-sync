package config

import (
	"context"
	"database/sql"
	"trakt-sync/internal/ctxutils"
)

func InitConfigTable(ctx *context.Context) *ConfigEntity {
	// Retrieve the database connection from the context
	db, ok := (*ctx).Value(ctxutils.ContextDbKey).(*sql.DB)
	if !ok || db == nil {
		panic("Database connection not found in context")
	}

	// Create the config table if it does not exist
	createTableQuery := `CREATE TABLE IF NOT EXISTS config (
		type TEXT PRIMARY KEY,
		data TEXT
	)`
	if _, err := db.Exec(createTableQuery); err != nil {
		panic("Failed to create config table: " + err.Error())
	}

	cfg, err := ReadConfig(ctx)
	if err != nil {
		panic("Failed to read config: " + err.Error())
	}

	if cfg.Trakt == nil && cfg.Emby == nil && cfg.Plex == nil && cfg.Jellyfin == nil {
		// Initialize with default values if no cfg exists
		cfg := ConfigEntity{
			Trakt: &TraktConfig{
				ClientID: "eb4ede9a384157e9aa60aad8c72c36c0485215659c82ad7b1fe965359a55caf4",
				// ClientSecret: uuid.New().String(),
				// AccessToken:  uuid.New().String(),
				// RefreshToken: uuid.New().String(),
				// Code:         uuid.New().String(),
				RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
			},
			Emby: &EmbyConfig{
				// UserID:  uuid.New().String(),
				// APIKey:  uuid.New().String(),
				BaseURL: "http://localhost:8096",
			},
			Plex: &PlexConfig{
				BaseURL: "http://localhost:32400",
				// APIKey:  "your-plex-api-key",
			},
			Jellyfin: &JellyfinConfig{
				BaseURL: "http://localhost:8920",
				// APIKey:  "your-jellyfin-api-key",
			},
			Cronjob: "0 0 * * *", // Default cronjob to run daily at midnight
		}
		UpsertConfig(ctx, &cfg)
	}

	return cfg
}
