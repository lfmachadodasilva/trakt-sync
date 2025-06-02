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

	config.InitConfigTable(&ctx)

	http.HandleFunc("/config", HandleConfig(&ctx))
	http.HandleFunc("/emby/", HandleEmby(&ctx))
	http.HandleFunc("/trakt/", HandleTrakt(&ctx))

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("🚀 Starting server on port 3000...")
	http.ListenAndServe(":3000", nil)

	// Close the database connection
	db := database.GetAndConnect(&ctx)
	defer db.Close()
}
