package main

import (
	"context"
	"fmt"
	"net/http"
	"trakt-sync/internal/config"
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
	ctx = context.WithValue(ctx, "db", db)

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

	fmt.Println("🚀 Starting server on port 3000...")
	http.ListenAndServe(":3000", nil)

	// Close the database connection
	defer db.Close()
}
