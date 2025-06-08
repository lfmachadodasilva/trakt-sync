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

func MarkItemAsWatched(ctx *context.Context, cfg *config.ConfigEntity, request *MarkAsWatchedRequest) error {

	preUrl := "%s/sync/history"
	url := fmt.Sprintf(preUrl, TraktApiUrl)

	_, err := utils.HttpPost[MarkAsWatchedRequest, struct{}](
		utils.RequestParams{
			URL:        url,
			Config:     cfg,
			AddHeaders: addTraktHeaders,
			Context:    ctx,
		},
		request,
	)
	if err != nil {
		return fmt.Errorf("failed to mark trakt item as watched: %w", err)
	}

	return nil
}
