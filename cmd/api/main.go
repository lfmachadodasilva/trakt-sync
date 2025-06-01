package main

import (
	"fmt"
	"net/http"
	"trakt-sync/internal/database"
)

func main() {
	database.InitDatabase() // Initialize the database and create the config table if it does not exist

	http.HandleFunc("/config", HandleConfig())
	http.HandleFunc("/emby/", HandleEmby())

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("🚀 Starting server on port 3000...")
	http.ListenAndServe(":3000", nil)

	// Close the database connection
	db := database.GetAndConnect()
	defer db.Close()
}
