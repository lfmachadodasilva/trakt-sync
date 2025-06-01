package trakt

import (
	"fmt"
	"time"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

type TraktWatchedItem struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	Ids   struct {
		Trakt int    `json:"trakt"`
		Slug  string `json:"slug"`
		Imdb  string `json:"imdb"`
	} `json:"ids"`
}

type TraktWatchedEpisode struct {
	Number        int       `json:"number"`
	LastWatchedAt time.Time `json:"last_watched_at"`
}

type TraktWatchedSeason struct {
	Number   int                   `json:"number"`
	Episodes []TraktWatchedEpisode `json:"episodes"`
}

type TraktWatchedResponse struct {
	LastWatchedAt time.Time             `json:"last_watched_at"`
	LastUpdatedAt time.Time             `json:"last_updated_at"`
	Movie         *TraktWatchedItem     `json:"movie"`
	Show          *TraktWatchedItem     `json:"show"`
	Seasons       *[]TraktWatchedSeason `json:"seasons"`
}

type TraktWatched struct {
	Movies []TraktWatchedResponse
	Shows  []TraktWatchedResponse
}

// FetchTraktWatched fetches the watched movies and shows from Trakt using the provided config.Config
// Documentation: https://trakt.docs.apiary.io/#reference/sync/get-watched/get-watched
func FetchTraktWatched(config *config.ConfigEntity) (*TraktWatched, error) {

	movies, err := getWatchedGeneric(config, "movies")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch watched movies: %w", err)
	}
	shows, err := getWatchedGeneric(config, "shows")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch watched shows: %w", err)
	}

	return &TraktWatched{
		Movies: *movies,
		Shows:  *shows,
	}, nil
}

func getWatchedGeneric(config *config.ConfigEntity, mediaType string) (*[]TraktWatchedResponse, error) {

	// Validate the itemType parameter
	if mediaType != "movies" && mediaType != "shows" {
		return nil, fmt.Errorf("invalid itemType: %s. Must be 'movies' or 'shows'", mediaType)
	}

	url := TraktApiUrl + "/sync/watched/" + mediaType

	response, err := utils.Get[[]TraktWatchedResponse](url, config, addTraktHeaders)
	if err != nil {
		return nil, err
	}

	return response, nil
}
