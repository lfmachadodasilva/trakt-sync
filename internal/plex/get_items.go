package plex

import (
	"context"
	"trakt-sync/internal/config"
)

type PlexMovieResponse struct {
	Id      string
	Title   string
	ImdbId  string
	Watched bool
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
	return nil, nil
}
