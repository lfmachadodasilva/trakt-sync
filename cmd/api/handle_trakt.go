package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"trakt-sync/internal/config"
	"trakt-sync/internal/trakt"
)

func HandleTrakt(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the sub-path after /trakt/
		subPath := r.URL.Path[len("/trakt"):]

		switch subPath {
		case "/code":
			// Handle the base /trakt endpoint
			switch r.Method {
			case http.MethodGet:
				HandleTraktCode(ctx, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case "/auth":
			// Handle the base /trakt endpoint
			switch r.Method {
			case http.MethodPost:
				HandleTraktAuth(ctx, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case "/auth/refresh":
			// Handle the base /trakt endpoint
			switch r.Method {
			case http.MethodPost:
				HandleTraktAuthRefresh(ctx, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Handle other sub-paths under /trakt/
			http.Error(w, "Endpoint not found", http.StatusNotFound)
		}
	}
}

func HandleTraktCode(ctx *context.Context, w http.ResponseWriter, r *http.Request) {

	cfg, err := config.ReadConfig(ctx)
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusInternalServerError)
		return
	}

	if cfg.Trakt.ClientID == "" || cfg.Trakt.RedirectURL == "" {
		http.Error(w, "Trakt ClientID or RedirectURL is not set", http.StatusBadRequest)
		return
	}

	preUrl := "%s/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s"
	url := fmt.Sprintf(preUrl, trakt.TraktApiUrl, cfg.Trakt.ClientID, cfg.Trakt.RedirectURL)

	w.Write([]byte(url))
}

func HandleTraktAuth(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	cfg, err := config.ReadConfig(ctx)
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusInternalServerError)
		return
	}

	// Define a struct to hold the expected JSON body
	var requestBody struct {
		Code string `json:"code"`
	}

	// Decode the JSON body into the struct
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = trakt.Auth(ctx, cfg, requestBody.Code)
	if err != nil {
		http.Error(w, "Failed to fetch Trakt auth: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleTraktAuthRefresh(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	cfg, err := config.ReadConfig(ctx)
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusInternalServerError)
		return
	}

	err = trakt.AuthRefreshAccessToken(ctx, cfg)
	if err != nil {
		http.Error(w, "Failed to fetch Trakt auth: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
