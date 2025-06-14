package trakt

import (
	"context"
	"fmt"
	"time"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

type MarkAsWatchedIds struct {
	Imdb string `json:"imdb"`
}

type MarkAsWatchedMovieRequest struct {
	Ids       MarkAsWatchedIds `json:"ids"`
	WatchedAt time.Time        `json:"watched_at"`
}

type MarkAsWatchedShowRequest struct {
	Ids     MarkAsWatchedIds              `json:"ids"`
	Seasons []MarkAsWatchedSeasonsRequest `json:"seasons,omitempty"`
}

type MarkAsWatchedSeasonsRequest struct {
	WatchedAt time.Time               `json:"watched_at"`
	Number    int16                   `json:"number"`
	Episodes  []MarkAsWatchedEpisodes `json:"episodes,omitempty"`
}

type MarkAsWatchedEpisodes struct {
	Number    int16     `json:"number"`
	WatchedAt time.Time `json:"watched_at"`
}

type MarkAsWatchedRequest struct {
	Movies []MarkAsWatchedMovieRequest `json:"movies,omitempty"`
	Shows  []MarkAsWatchedShowRequest  `json:"shows,omitempty"`
}

type MarkAsWatchedResponse struct {
	Added struct {
		Movies   int16 `json:"movies,omitempty"`
		Episodes int16 `json:"episodes,omitempty"`
	} `json:"added,omitempty"`
}

func MarkItemAsWatched(ctx *context.Context, cfg *config.ConfigEntity, request *MarkAsWatchedRequest, isRetry bool) error {

	preUrl := "%s/sync/history"
	url := fmt.Sprintf(preUrl, TraktApiUrl)

	response, err := utils.HttpPost[MarkAsWatchedRequest, MarkAsWatchedResponse](
		utils.RequestParams{
			URL:        url,
			Config:     cfg,
			AddHeaders: addTraktHeaders,
			Context:    ctx,
		},
		request,
	)
	if err != nil {
		if !isRetry && utils.IsAuthError(err) {
			// If the error is an authentication error, refresh the access token and retry
			err = AuthRefreshAccessToken(ctx, cfg)
			if err != nil {
				return err
			}
			// cfg, err = config.GetConfig(ctx)
			// Retry the request after refreshing the access token
			return MarkItemAsWatched(ctx, cfg, request, true)
		}
	} else if response.Added.Movies == 0 && response.Added.Episodes == 0 {
		return fmt.Errorf("no items were added to the watch history, check if the request is valid")
	}

	return err
}
