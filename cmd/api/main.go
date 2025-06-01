package main

import (
	"fmt"
	"net/http"
	"trakt-sync/internal/config"
	"trakt-sync/internal/database"
	"trakt-sync/internal/emby"
	"trakt-sync/internal/trakt"
)

func main() {
	config.InitConfigTable() // Initialize the database and create the config table if it does not exist

	http.HandleFunc("/config", HandleConfig())
	http.HandleFunc("/emby/", HandleEmby())
	http.HandleFunc("/trakt/", HandleTrakt())

	// DEBUG: Fetch Emby items for testing purposes
	config, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	users, err := emby.FetchEmbyUsers(&config) // Fetch Emby users for testing purposes
	if err != nil {
		fmt.Println("Error fetching users:", err)
		return
	}
	for _, user := range users {
		if user.Name == "luizfelipe" {
			config.Emby.UserID = user.Id // Set the UserID in the config to the target user's ID
			break
		}
	}
	movies, err := emby.FetchEmbyItems(&config, "Movie") // Fetch Emby items of type Movie
	if err != nil {
		fmt.Println("Error fetching movies:", err)
		return
	}
	fmt.Println("Fetched Movies:", movies)
	series, err := emby.FetchEmbyItems(&config, "Series") // Fetch Emby items of type Series
	if err != nil {
		fmt.Println("Error fetching series:", err)
		return
	}
	fmt.Println("Fetched Series:", series)
	items, err := emby.FetchEmbyItemsFull(&config) // Fetch all Emby items of type Movie and Series
	if err != nil {
		fmt.Println("Error fetching all Emby items:", err)
		return
	}
	fmt.Println("Fetched All Emby Items:", items)
	watched, err := trakt.FetchTraktWatched(&config) // Fetch watched items from Emby
	if err != nil {
		fmt.Println("Error fetching watched items:", err)
		return
	}
	fmt.Println("Fetched Watched Items:", watched)
	// DEBUG END

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("🚀 Starting server on port 3000...")
	http.ListenAndServe(":3000", nil)

	// Close the database connection
	db := database.GetAndConnect()
	defer db.Close()
}
