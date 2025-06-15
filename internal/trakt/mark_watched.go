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
	Number   int16                   `json:"number"`
	Episodes []MarkAsWatchedEpisodes `json:"episodes,omitempty"`
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
	} `json:"added"`
}

type MarkAsWatchedMap struct {
	Movies map[string]time.Time
	Shows  map[string]map[int16]map[int16]time.Time
}

func MarkItemAsWatched(ctx *context.Context, cfg *config.ConfigEntity, request *MarkAsWatchedMap, isRetry bool) error {

	markAsWatchedRequest := request.ToMarkAsWatchedRequest()

	preUrl := "%s/sync/history"
	url := fmt.Sprintf(preUrl, TraktApiUrl)

	response, err := utils.HttpPost[MarkAsWatchedRequest, MarkAsWatchedResponse](
		utils.RequestParams{
			URL:        url,
			Config:     cfg,
			AddHeaders: addTraktHeaders,
			Context:    ctx,
		},
		markAsWatchedRequest,
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

func (request *MarkAsWatchedMap) AppendMovie(imdbId string, watchedAt time.Time) error {

	if request.Movies == nil {
		request.Movies = make(map[string]time.Time)
	}

	if _, exists := request.Movies[imdbId]; !exists {
		request.Movies[imdbId] = watchedAt
	}

	return nil
}

func (request *MarkAsWatchedMap) AppendTvShow(imdbId string, seasonNumber int16, episodeNumber int16, watchedAt time.Time) error {

	if request.Shows == nil {
		request.Shows = make(map[string]map[int16]map[int16]time.Time)
	}

	if _, exists := request.Shows[imdbId]; !exists {
		request.Shows[imdbId] = make(map[int16]map[int16]time.Time)
	}

	if _, exists := request.Shows[imdbId][seasonNumber]; !exists {
		request.Shows[imdbId][seasonNumber] = make(map[int16]time.Time)
	}

	if _, exists := request.Shows[imdbId][seasonNumber][episodeNumber]; !exists {
		request.Shows[imdbId][seasonNumber][episodeNumber] = watchedAt
	}

	return nil
}

func (request *MarkAsWatchedMap) ToMarkAsWatchedRequest() *MarkAsWatchedRequest {
	markAsWatchedRequest := MarkAsWatchedRequest{
		Movies: make([]MarkAsWatchedMovieRequest, 0),
		Shows:  make([]MarkAsWatchedShowRequest, 0),
	}

	if request.Movies != nil {
		for imdbId, watchedAt := range request.Movies {
			markAsWatchedRequest.Movies = append(markAsWatchedRequest.Movies, MarkAsWatchedMovieRequest{
				Ids: MarkAsWatchedIds{
					Imdb: imdbId,
				},
				WatchedAt: watchedAt,
			})
		}
	}

	if request.Shows != nil {
		for imdbId, seasons := range request.Shows {
			showRequest := MarkAsWatchedShowRequest{
				Ids: MarkAsWatchedIds{
					Imdb: imdbId,
				},
			}
			for season, episodes := range seasons {
				seasonRequest := MarkAsWatchedSeasonsRequest{
					Number: season,
				}
				for episode, watchedAt := range episodes {
					seasonRequest.Episodes = append(seasonRequest.Episodes, MarkAsWatchedEpisodes{
						Number:    episode,
						WatchedAt: watchedAt,
					})
				}
				showRequest.Seasons = append(showRequest.Seasons, seasonRequest)
			}
			markAsWatchedRequest.Shows = append(markAsWatchedRequest.Shows, showRequest)
		}
	}

	return &markAsWatchedRequest
}
