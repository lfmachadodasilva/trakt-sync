package trakt

import (
	"context"
	"fmt"
	"time"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

type TraktWatchedItem struct {
	Title string `json:"title"`
	Year  int16  `json:"year"`
	Ids   struct {
		Trakt int    `json:"trakt,omitempty"`
		Slug  string `json:"slug,omitempty"`
		Imdb  string `json:"imdb,omitempty"`
	} `json:"ids"`
}

type TraktWatchedEpisode struct {
	Number        int16     `json:"number,omitempty"`
	LastWatchedAt time.Time `json:"last_watched_at,omitempty"`
}

type TraktWatchedSeason struct {
	Number   int16                 `json:"number,omitempty"`
	Episodes []TraktWatchedEpisode `json:"episodes,omitempty"`
}

type TraktWatchedResponse struct {
	LastWatchedAt time.Time            `json:"last_watched_at"`
	LastUpdatedAt time.Time            `json:"last_updated_at"`
	Movie         TraktWatchedItem     `json:"movie,omitempty"`
	Show          TraktWatchedItem     `json:"show,omitempty"`
	Seasons       []TraktWatchedSeason `json:"seasons,omitempty"`
}

type TraktWatched struct {
	Movies []TraktWatchedResponse
	Shows  []TraktWatchedResponse
}

// GetWatched fetches the watched movies and shows from Trakt using the provided config.Config
// Documentation: https://trakt.docs.apiary.io/#reference/sync/get-watched/get-watched
func GetWatched(ctx *context.Context, cfg *config.ConfigEntity) (*TraktWatched, error) {

	movies, err := getWatchedGeneric(ctx, cfg, "movies")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch watched movies: %w", err)
	}
	shows, err := getWatchedGeneric(ctx, cfg, "shows")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch watched shows: %w", err)
	}

	return &TraktWatched{
		Movies: *movies,
		Shows:  *shows,
	}, nil
}

func getWatchedGeneric(ctx *context.Context, cfg *config.ConfigEntity, mediaType string) (*[]TraktWatchedResponse, error) {

	// Validate the itemType parameter
	if mediaType != "movies" && mediaType != "shows" {
		return nil, fmt.Errorf("invalid itemType: %s. Must be 'movies' or 'shows'", mediaType)
	}

	url := TraktApiUrl + "/sync/watched/" + mediaType

	response, err := utils.HttpGet[[]TraktWatchedResponse](
		utils.RequestParams{
			URL:        url,
			Config:     cfg,
			AddHeaders: addTraktHeaders,
			Context:    ctx,
		})
	if err != nil {
		return nil, err
	}

	return response, nil
}
