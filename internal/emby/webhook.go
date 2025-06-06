package emby

import (
	"context"
	"fmt"
	"strings"
	"time"
	"trakt-sync/internal/config"
	"trakt-sync/internal/trakt"
)

type EmbyWebhook struct {
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	Event string    `json:"event"`
	User  struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"user"`
	Item struct {
		Name              string            `json:"name"`
		Id                string            `json:"id"`
		ProviderIds       map[string]string `json:"provider_ids"`
		IndexNumber       *int              `json:"index_number,omitempty"`
		ParentIndexNumber *int              `json:"parent_index_number,omitempty"`
		Type              string            `json:"type"`
	} `json:"item"`
	Server map[string]string `json:"server"`
}

func (webhook *EmbyWebhook) GetImdbId() (string, error) {

	imdbId := webhook.Item.ProviderIds["imdb"]
	if imdbId == "" {
		imdbId = webhook.Item.ProviderIds["IMDB"]
	}

	if imdbId == "" {
		return "", fmt.Errorf("no IMDb ID found for item: %s", webhook.Item.Name)
	}
	return imdbId, nil
}

func ProcessEmbyWebhook(ctx *context.Context, cfg *config.ConfigEntity, webhook *EmbyWebhook) error {
	fmt.Printf("Received Emby webhook: %s at %s by user %s - %s - %s\n", webhook.Item.Name, webhook.Date, webhook.User.Name, webhook.Event, webhook.Title)

	if webhook.User.Id != cfg.Emby.UserID {
		fmt.Println("Ignoring webhook for user:", webhook.User.Name, "with ID:", webhook.User.Id)
		return nil
	}

	shouldProcess := strings.Contains(webhook.Title, "has finished playing")

	if !shouldProcess {
		fmt.Println("Ignoring webhook event:", webhook.Event, "with title:", webhook.Title)
		return nil
	}

	if err := processEmbyMovie(ctx, cfg, webhook); err != nil {
		return fmt.Errorf("failed to process Emby movie: %w", err)
	}
	if err := processEmbySeries(ctx, cfg, webhook); err != nil {
		return fmt.Errorf("failed to process Emby series: %w", err)
	}

	return nil
}

func processEmbyMovie(ctx *context.Context, cfg *config.ConfigEntity, webhook *EmbyWebhook) error {
	if webhook.Item.Type != "Movie" {
		return nil
	}

	imdbId, err := webhook.GetImdbId()
	if err != nil {
		return fmt.Errorf("failed to get IMDb ID for Emby movie: %w", err)
	}

	fmt.Printf("Processing Emby movie: %s with IMDb ID: %s\n", webhook.Item.Name, imdbId)
	traktRequest := &trakt.MarkAsWatchedRequest{
		Movies: []trakt.MarkAsWatchedMovieRequest{
			{
				Ids: trakt.MarkAsWatchedIds{
					Imdb: imdbId,
				},
				WatchedAt: webhook.Date,
			},
		},
	}
	if err := trakt.MarkItemAsWatched(ctx, cfg, traktRequest); err != nil {
		return fmt.Errorf("failed to mark Emby movie as watched in Trakt: %w", err)
	}
	fmt.Printf("Marked Emby movie: %s as watched in Trakt with IMDB: %s\n", webhook.Item.Name, imdbId)

	return nil
}

func processEmbySeries(ctx *context.Context, cfg *config.ConfigEntity, webhook *EmbyWebhook) error {
	if webhook.Item.Type != "Series" {
		return nil
	}

	imdbId, err := webhook.GetImdbId()
	if err != nil {
		return fmt.Errorf("failed to get IMDb ID for Emby series: %w", err)
	}

	fmt.Printf("Processing Emby series: %s with IMDb ID: %s\n", webhook.Item.Name, imdbId)
	traktRequest := &trakt.MarkAsWatchedRequest{
		Shows: []trakt.MarkAsWatchedShowRequest{
			{
				Ids: trakt.MarkAsWatchedIds{
					Imdb: imdbId,
				},
				Seasons: []trakt.MarkAsWatchedSeasonsRequest{
					{
						WatchedAt: webhook.Date,
						Number:    int16(*webhook.Item.ParentIndexNumber),
						Episodes: []trakt.MarkAsWatchedEpisodes{
							{
								Number:    int16(*webhook.Item.IndexNumber),
								WatchedAt: webhook.Date,
							},
						},
					},
				},
			},
		},
	}
	if err := trakt.MarkItemAsWatched(ctx, cfg, traktRequest); err != nil {
		return fmt.Errorf("failed to mark Emby series as watched in Trakt: %w", err)
	}
	fmt.Printf("Marked Emby series: %s as watched in Trakt with IMDB: %s\n", webhook.Item.Name, imdbId)

	return nil
}
