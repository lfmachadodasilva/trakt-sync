package trakt

import (
	"fmt"
	"time"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

type MarkAsWatchedRequest struct {
	Movies *[]struct {
		Ids struct {
			Imdb string `json:"imdb"`
		} `json:"ids"`
		WatchedAt time.Time `json:"watched_at"`
	} `json:"movies,omitempty"`

	Shows *[]struct {
		Ids struct {
			Imdb string `json:"imdb"`
		} `json:"ids"`
		Seasons *[]struct {
			WatchedAt time.Time `json:"watched_at"`
			Number    int       `json:"number"`
			Episodes  *[]struct {
				Number    int       `json:"number"`
				WatchedAt time.Time `json:"watched_at"`
			} `json:"episodes,omitempty"`
		} `json:"seasons,omitempty"`
	} `json:"shows,omitempty"`
}

func MarkItemAsWatched(cfg *config.ConfigEntity, request *MarkAsWatchedRequest) error {

	preUrl := "%s/sync/history"
	url := fmt.Sprintf(preUrl, TraktApiUrl)

	_, err := utils.HttpPost[MarkAsWatchedRequest, struct{}](
		utils.RequestParams{
			URL:        url,
			Config:     cfg,
			AddHeaders: addTraktHeaders,
		},
		request,
	)
	if err != nil {
		return fmt.Errorf("failed to mark item as watched: %w", err)
	}

	return nil
}
