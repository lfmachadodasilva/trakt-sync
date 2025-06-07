package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"trakt-sync/internal/config"
	"trakt-sync/internal/ctxutils"
	"trakt-sync/internal/database"
)

func main() {

	ctx := context.Background()

	db := database.GetAndConnect(&ctx)
	if db == nil {
		fmt.Println("Failed to connect to the database")
		return
	}

	// Add the database connection to the context
	ctx = context.WithValue(ctx, ctxutils.ContextDbKey, db)

	// Initialize the configuration table
	cfg := config.InitConfigTable(&ctx)

	if cfg.Cronjob != "" {
		cronManager := config.NewCronManager()
		defer cronManager.Stop()

		ctx = context.WithValue(ctx, "cron", cronManager)

		cronManager.Start(context.Background(), cfg.Cronjob, func() {
			fmt.Println("Running sync job every minute...")

			cfg, err := config.ReadConfig(&ctx)
			if err != nil {
				fmt.Println("Failed to read config:", err)
				return
			}

			fmt.Printf("Running config: %s\n", cfg.Cronjob)

			// Here you would call the sync function to perform the sync
			// For example: sync.SyncAll(context.Background())
		})
	}

	http.HandleFunc("/config", HandleConfig(&ctx))
	http.HandleFunc("/emby/", HandleEmby(&ctx))
	http.HandleFunc("/trakt/", HandleTrakt(&ctx))
	http.HandleFunc("/sync", HandleSync(&ctx))

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Determine the port based on the environment
	port := "4000"
	if os.Getenv("NODE_ENV") == "production" {
		port = "3000"
	}

	fmt.Printf("🚀 Starting server on port %s...\n", port)
	http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux))

	// Close the database connection
	defer db.Close()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (you can restrict this to specific domains)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
