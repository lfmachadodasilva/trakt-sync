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
	Id          string
	Title       string
	Season      int16
	Episode     int16
	Watched     bool
	WatchedTime time.Time
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

	responseShows, err := sdk.Library.GetLibrarySectionsAll(*ctx, operations.GetLibrarySectionsAllRequest{
		SectionKey:           2,
		Type:                 operations.GetLibrarySectionsAllQueryParamTypeEpisode,
		IncludeMeta:          operations.GetLibrarySectionsAllQueryParamIncludeMetaEnable.ToPointer(),
		IncludeGuids:         operations.QueryParamIncludeGuidsEnable.ToPointer(),
		IncludeAdvanced:      operations.IncludeAdvancedEnable.ToPointer(),
		IncludeCollections:   operations.QueryParamIncludeCollectionsEnable.ToPointer(),
		IncludeExternalMedia: operations.QueryParamIncludeExternalMediaEnable.ToPointer(),
	})
	if err != nil {
		log.Fatal(err)
	}

	var tvShows []PlexTvshowResponse

	for _, show := range responseShows.Object.MediaContainer.Metadata {
		var imdbId string
		for _, guid := range show.Guids {
			if strings.HasPrefix(guid.ID, "imdb://") {
				// extract the imdb id from the guid
				imdbId = strings.TrimPrefix(guid.ID, "imdb://")
				break
			}
		}
		if imdbId == "" {
			fmt.Printf("no imdb id found for plex tv show: %s\n", show.Title)
			continue
		}
		// var episodes []PlexEpisodeResponse
		// for _, episode := range show.Children {
		// 	var watched bool = false
		// 	if episode.ViewCount != nil && *episode.ViewCount > 0 {
		// 		watched = true
		// 	}
		// 	var watchedTime time.Time
		// 	if episode.LastViewedAt != nil {
		// 		watchedTime = time.Unix(int64(*episode.LastViewedAt), 0)
		// 	} else {
		// 		watchedTime = time.Time{}
		// 	}
		// 	episodes = append(episodes, PlexEpisodeResponse{
		// 		Id:          episode.RatingKey,
		// 		Title:       episode.Title,
		// 		Season:      episode.ParentIndex,
		// 		Episode:     episode.Index,
		// 		Watched:     watched,
		// 		WatchedTime: watchedTime,
		// 	})
		// }
		// tvShows = append(tvShows, PlexTvshowResponse{
		// 	Id:       show.RatingKey,
		// 	Title:    show.Title,
		// 	ImdbId:   imdbId,
		// 	Episodes: episodes,
		// })
	}

	return tvShows, nil
}
