package main

import (
	"encoding/json"
	"net/http"
	"trakt-sync/internal/database"
	"trakt-sync/internal/emby"
)

func HandleEmby() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the sub-path after /emby/
		subPath := r.URL.Path[len("/emby"):]

		switch subPath {
		case "/users":
			// Handle the base /emby endpoint
			switch r.Method {
			case http.MethodGet:
				HandleEmbyUsers(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Handle other sub-paths under /emby/
			http.Error(w, "Endpoint not found", http.StatusNotFound)
		}
	}
}

func HandleEmbyUsers(w http.ResponseWriter, r *http.Request) {

	config, err := database.ReadConfig()
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusInternalServerError)
		return
	}
	users, err := emby.FetchEmbyUsers(&config)
	if err != nil {
		http.Error(w, "Failed to fetch Emby users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the fetched users into JSON
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to marshal users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	w.WriteHeader(http.StatusOK)
}
