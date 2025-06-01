package main

import (
	"encoding/json"
	"net/http"
	"trakt-sync/internal/database"
	"trakt-sync/internal/models"
)

func HandleConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the sub-path after /config/
		subPath := r.URL.Path[len("/config"):]

		switch subPath {
		case "":
			// Handle the base /config endpoint
			switch r.Method {
			case http.MethodGet:
				handleGetConfig(w)
			case http.MethodPatch:
				handlePatchConfig(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Handle other sub-paths under /config/
			http.Error(w, "Endpoint not found", http.StatusNotFound)
		}
	}
}

func handleGetConfig(w http.ResponseWriter) {
	config, err := database.ReadConfig()
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusInternalServerError)
		return
	}

	// Marshal the fetched configurations into JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		http.Error(w, "Failed to marshal configs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func handlePatchConfig(w http.ResponseWriter, r *http.Request) {
	var request models.Config
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := database.UpsertConfig(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
