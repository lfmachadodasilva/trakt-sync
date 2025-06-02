package main

import (
	"context"
	"encoding/json"
	"net/http"
	"trakt-sync/internal/config"
)

func HandleConfig(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the sub-path after /config/
		subPath := r.URL.Path[len("/config"):]

		switch subPath {
		case "":
			// Handle the base /config endpoint
			switch r.Method {
			case http.MethodGet:
				handleGetConfig(ctx, w)
			case http.MethodPatch:
				handlePatchConfig(ctx, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Handle other sub-paths under /config/
			http.Error(w, "Endpoint not found", http.StatusNotFound)
		}
	}
}

func handleGetConfig(ctx *context.Context, w http.ResponseWriter) {
	cfg, err := config.ReadConfig(ctx)
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusInternalServerError)
		return
	}

	// Marshal the fetched configurations into JSON
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		http.Error(w, "Failed to marshal configs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func handlePatchConfig(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	var request config.ConfigEntity
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := config.UpsertConfig(ctx, &request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
