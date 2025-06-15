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
	Movies map[string]*TraktWatchedItem
	Shows  map[string]map[int16]map[int16]*TraktWatchedResponse
}

// GetWatched fetches the watched movies and shows from Trakt using the provided config.Config
// Documentation: https://trakt.docs.apiary.io/#reference/sync/get-watched/get-watched
func GetWatched(ctx *context.Context, cfg *config.ConfigEntity) (TraktWatched, error) {

	movies, err := getWatchedGeneric(ctx, cfg, "movies", false)
	if err != nil {
		return TraktWatched{}, fmt.Errorf("failed to fetch watched movies: %w", err)
	}
	shows, err := getWatchedGeneric(ctx, cfg, "shows", false)
	if err != nil {
		return TraktWatched{}, fmt.Errorf("failed to fetch watched shows: %w", err)
	}

	return TraktWatched{
		Movies: CreateTraktMovieMap(movies),
		Shows:  CreateTraktShowMap(shows),
	}, nil
}

func getWatchedGeneric(ctx *context.Context, cfg *config.ConfigEntity, mediaType string, isRetry bool) (*[]TraktWatchedResponse, error) {

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
		if !isRetry && utils.IsAuthError(err) {
			// If the error is an authentication error, refresh the access token and retry
			err = AuthRefreshAccessToken(ctx, cfg)
			if err != nil {
				return nil, fmt.Errorf("failed to refresh access token: %w", err)
			}
			// cfg, err = config.GetConfig(ctx)
			// Retry the request after refreshing the access token
			return getWatchedGeneric(ctx, cfg, mediaType, true)
		}

		return nil, err
	}

	return response, nil
}

// CreateTraktMovieMap creates a map with IMDb IDs as keys and TraktWatchedItem as values
func CreateTraktMovieMap(traktItems *[]TraktWatchedResponse) map[string]*TraktWatchedItem {
	imdbMap := make(map[string]*TraktWatchedItem)

	// Iterate over TraktItems (movies) and populate the map
	for _, movie := range *traktItems {
		if movie.Movie.Ids.Imdb != "" {
			imdbMap[movie.Movie.Ids.Imdb] = &movie.Movie
		}
	}

	return imdbMap
}

// CreateTraktShowMap creates a nested map with IMDb IDs, season numbers, and episode numbers as keys and TraktWatchedResponse as values
func CreateTraktShowMap(traktItems *[]TraktWatchedResponse) map[string]map[int16]map[int16]*TraktWatchedResponse {
	traktImdbMap := make(map[string]map[int16]map[int16]*TraktWatchedResponse)

	for _, show := range *traktItems {
		if show.Show.Ids.Imdb != "" {
			imdbId := show.Show.Ids.Imdb
			if _, exists := traktImdbMap[imdbId]; !exists {
				traktImdbMap[imdbId] = make(map[int16]map[int16]*TraktWatchedResponse)
			}

			for _, season := range show.Seasons {
				seasonNumber := season.Number
				if _, exists := traktImdbMap[imdbId][seasonNumber]; !exists {
					traktImdbMap[imdbId][seasonNumber] = make(map[int16]*TraktWatchedResponse)
				}

				for _, episode := range season.Episodes {
					episodeNumber := episode.Number
					traktImdbMap[imdbId][seasonNumber][episodeNumber] = &show
				}
			}
		}
	}

	return traktImdbMap
}
