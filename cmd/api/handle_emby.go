package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"trakt-sync/internal/config"
	"trakt-sync/internal/emby"
	"trakt-sync/internal/utils"
)

func HandleEmby(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the sub-path after /emby/
		subPath := r.URL.Path[len("/emby"):]

		switch subPath {
		case "/users":
			// Handle the base /emby/users endpoint
			switch r.Method {
			case http.MethodGet:
				HandleEmbyUsers(ctx, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case "/webhooks":
			// Handle the base /emby/webhooks endpoint
			switch r.Method {
			case http.MethodPost:
				HandleEmbyWebhooks(ctx, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Handle other sub-paths under /emby/
			http.Error(w, "Endpoint not found", http.StatusNotFound)
		}
	}
}

func HandleEmbyUsers(ctx *context.Context, w http.ResponseWriter, r *http.Request) {

	cfg, err := config.ReadConfig(ctx)
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusInternalServerError)
		fmt.Println("Failed to read configs:", err)
		return
	}
	usr, err := emby.FetchEmbyUsers(cfg)
	if err != nil {
		http.Error(w, "Failed to fetch Emby users: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Failed to fetch Emby users:", err)
		return
	}

	// Marshal the fetched users into JSON
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(usr)
	if err != nil {
		http.Error(w, "Failed to marshal users", http.StatusInternalServerError)
		fmt.Println("Failed to marshal users:", err)
		return
	}

	w.Write(jsonData)
}

func HandleEmbyWebhooks(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	cfg, err := config.ReadConfig(ctx)
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusInternalServerError)
		return
	}

	webhook, err := utils.SerializeBody[emby.EmbyWebhook](r.Body)
	if err != nil {
		http.Error(w, "Failed to parse webhook: "+err.Error(), http.StatusBadRequest)
		fmt.Println("Failed to parse webhook:", err)
		return
	}

	err = emby.ProcessEmbyWebhook(ctx, cfg, webhook)
	if err != nil {
		http.Error(w, "Failed to process webhook: "+err.Error(), http.StatusBadGateway)
		fmt.Println("Failed to process webhook:", err)
		return
	}
}
