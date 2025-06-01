package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trakt-sync/internal/database"
	"trakt-sync/internal/models"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("API application started")

	config := models.Config{
		Trakt: &models.TraktConfig{
			ClientID:     uuid.New().String(),
			ClientSecret: uuid.New().String(),
			AccessToken:  uuid.New().String(),
			RefreshToken: uuid.New().String(),
			Code:         uuid.New().String(),
			RedirectURL:  "https://example.com/callback",
		},
		Emby: &models.EmbyConfig{
			UserID:  uuid.New().String(),
			APIKey:  uuid.New().String(),
			BaseURL: "https://emby.example.com",
		},
		Plex: &models.PlexConfig{
			UserID: uuid.New().String(),
		},
		Jellyfin: &models.JellyfinConfig{
			UserID: uuid.New().String(),
		},
	}

	db := database.Connect()
	configs := map[string]interface{}{
		"trakt":    config.Trakt,
		"emby":     config.Emby,
		"plex":     config.Plex,
		"jellyfin": config.Jellyfin,
	}

	for key, value := range configs {
		jsonData, err := json.Marshal(value)
		if err != nil {
			fmt.Printf("Failed to marshal %s config: %v\n", key, err)
			continue
		}
		database.UpsertConfig(db, key, string(jsonData))
	}

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			configs := map[string]interface{}{
				"trakt":    config.Trakt,
				"emby":     config.Emby,
				"plex":     config.Plex,
				"jellyfin": config.Jellyfin,
			}
			jsonData, err := json.Marshal(configs)
			if err != nil {
				http.Error(w, "Failed to marshal configs", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Starting server on port 3000...")
	http.ListenAndServe(":3000", nil)
}
