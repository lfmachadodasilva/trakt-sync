package plex

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"trakt-sync/internal/config"
	"trakt-sync/internal/ctxutils"

	"github.com/LukeHagar/plexgo"
	"github.com/LukeHagar/plexgo/models/operations"
)

type PlexMovieResponse struct {
	Id          string
	Title       string
	ImdbId      string
	Watched     bool
	WatchedTime time.Time
}

type PlexTvshowResponse struct {
	Id       string
	Title    string
	ImdbId   string
	Episodes []PlexEpisodeResponse
}

type PlexEpisodeResponse struct {
	Id      string
	Title   string
	Season  int16
	Episode int16
	Watched bool
}

type PlexItemsResponse struct {
	Movies  []PlexMovieResponse
	TvShows []PlexTvshowResponse
}

func GetAllItems(ctx *context.Context, cfg *config.ConfigEntity) (*PlexItemsResponse, error) {

	plexSdk, ok := (*ctx).Value(ctxutils.ContextPlexSdkKey).(*plexgo.PlexAPI)
	if !ok {
		return nil, fmt.Errorf("failed to get plex sdk from context")
	}

	movies, err := GetMovies(ctx, cfg, plexSdk)
	if err != nil {
		return nil, err
	}
	tvShows, err := GetTvShows(ctx, cfg, plexSdk)
	if err != nil {
		return nil, err
	}

	return &PlexItemsResponse{
		Movies:  movies,
		TvShows: tvShows,
	}, nil
}

func GetMovies(ctx *context.Context, cfg *config.ConfigEntity, sdk *plexgo.PlexAPI) ([]PlexMovieResponse, error) {

	responseMovies, err := sdk.Library.GetLibrarySectionsAll(*ctx, operations.GetLibrarySectionsAllRequest{
		SectionKey:           1,
		Type:                 operations.GetLibrarySectionsAllQueryParamTypeMovie,
		IncludeMeta:          operations.GetLibrarySectionsAllQueryParamIncludeMetaEnable.ToPointer(),
		IncludeGuids:         operations.QueryParamIncludeGuidsEnable.ToPointer(),
		IncludeAdvanced:      operations.IncludeAdvancedEnable.ToPointer(),
		IncludeCollections:   operations.QueryParamIncludeCollectionsEnable.ToPointer(),
		IncludeExternalMedia: operations.QueryParamIncludeExternalMediaEnable.ToPointer(),
	})
	if err != nil {
		log.Fatal(err)
	}
	if responseMovies.Object != nil {
		// handle response
	}

	// This function should implement the logic to fetch movies from Plex
	var movies []PlexMovieResponse
	for _, movie := range responseMovies.Object.MediaContainer.Metadata {
		var imdbId string
		for _, guid := range movie.Guids {
			if strings.HasPrefix(guid.ID, "imdb://") {
				// extract the imdb id from the guid
				imdbId = strings.TrimPrefix(guid.ID, "imdb://")
				break
			}
		}

		if imdbId == "" {
			fmt.Printf("no imdb id found for plex movie: %s\n", movie.Title)
			continue
		}

		var watched bool = false
		if movie.ViewCount != nil && *movie.ViewCount > 0 {
			watched = true
		}

		var watchedTime time.Time
		if movie.LastViewedAt != nil {
			watchedTime = time.Unix(int64(*movie.LastViewedAt), 0)
		} else {
			watchedTime = time.Time{}
		}

		movies = append(movies, PlexMovieResponse{
			Id:          movie.RatingKey,
			Title:       movie.Title,
			ImdbId:      imdbId,
			Watched:     watched,
			WatchedTime: watchedTime,
		})
	}

	return movies, nil
}

func GetTvShows(ctx *context.Context, cfg *config.ConfigEntity, sdk *plexgo.PlexAPI) ([]PlexTvshowResponse, error) {
	var tvShows []PlexTvshowResponse
	// This function should implement the logic to fetch TV shows from Plex
	return tvShows, nil
}
