package main

import (
	"fmt"
	"net/http"
	"trakt-sync/internal/config"
	"trakt-sync/internal/trakt"
)

func HandleTrakt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the sub-path after /trakt/
		subPath := r.URL.Path[len("/trakt"):]

		switch subPath {
		case "/code":
			// Handle the base /trakt endpoint
			switch r.Method {
			case http.MethodPost:
				HandleTraktCode(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case "/auth":
			// Handle the base /trakt endpoint
			switch r.Method {
			case http.MethodPost:
				HandleTraktAuth(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Handle other sub-paths under /trakt/
			http.Error(w, "Endpoint not found", http.StatusNotFound)
		}
	}
}

func HandleTraktCode(w http.ResponseWriter, r *http.Request) {

	config, err := config.ReadConfig()
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusInternalServerError)
		return
	}

	if config.Trakt.ClientID == "" || config.Trakt.RedirectURL == "" {
		http.Error(w, "Trakt ClientID or RedirectURL is not set", http.StatusBadRequest)
		return
	}

	preUrl := "%s/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s"
	url := fmt.Sprintf(preUrl, trakt.TraktApiUrl, config.Trakt.ClientID, config.Trakt.RedirectURL)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(url))
	w.WriteHeader(http.StatusOK)
}

func HandleTraktAuth(w http.ResponseWriter, r *http.Request) {
	// This function will handle the /trakt/auth endpoint
	// Implementation goes here
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
