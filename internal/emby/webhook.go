package emby

import (
	"context"
	"fmt"
	"time"
	"trakt-sync/internal/config"
	"trakt-sync/internal/trakt"
)

type EmbyWebhook struct {
	Title string    `json:"Title"`
	Date  time.Time `json:"Date"`
	Event string    `json:"Event"`
	User  struct {
		Name string `json:"Name"`
		Id   string `json:"Id"`
	} `json:"User"`
	Item struct {
		Name              string            `json:"Name"`
		Id                string            `json:"Id"`
		ProviderIds       map[string]string `json:"ProviderIds"`
		IndexNumber       *int              `json:"IndexNumber,omitempty"`
		ParentIndexNumber *int              `json:"ParentIndexNumber,omitempty"`
		Type              string            `json:"Type"`
		RunTimeTicks      int64             `json:"RunTimeTicks"`
		SeriesName        string            `json:"SeriesName,omitempty"`
	} `json:"Item"`
	Server       map[string]string `json:"Server"`
	PlaybackInfo struct {
		PositionTicks int64 `json:"PositionTicks"`
	} `json:"PlaybackInfo"`
}

func (webhook *EmbyWebhook) GetImdbId() (string, error) {

	imdbId := webhook.Item.ProviderIds["Imdb"]
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

	var perc int = 0 // Calculate the percentage of playback completion
	if webhook.Item.RunTimeTicks > 0 && webhook.PlaybackInfo.PositionTicks > 0 {
		perc = int((webhook.PlaybackInfo.PositionTicks * 100) / webhook.Item.RunTimeTicks)
	}

	shouldProcess :=
		// if the event is playback.stop and the percentage is greater than 90
		((webhook.Event == "playback.stop") && perc > 90) ||
			// or if the event is item.markplayed
			webhook.Event == "item.markplayed"

	if !shouldProcess {
		fmt.Println("Ignoring webhook event:", webhook.Event, "with title:", webhook.Title)
		return nil
	}

	if err := processEmbyMovie(ctx, cfg, webhook); err != nil {
		return fmt.Errorf("failed to process Emby movie: %w", err)
	}
	if err := processEmbyEpisode(ctx, cfg, webhook); err != nil {
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
		return err
	}

	fmt.Printf("Processing Emby movie: %s with IMDb ID: %s\n", webhook.Item.Name, imdbId)

	traktRequest := trakt.MarkAsWatchedMap{}
	if err := traktRequest.AppendMovie(imdbId, webhook.Date); err != nil {
		return fmt.Errorf("failed to append Emby movie to trakt request: %w", err)
	}

	if err := trakt.MarkItemAsWatched(ctx, cfg, &traktRequest, false); err != nil {
		return err
	}
	fmt.Printf("Marked Emby movie: %s as watched in Trakt with IMDB: %s\n", webhook.Item.Name, imdbId)

	return nil
}

func processEmbyEpisode(ctx *context.Context, cfg *config.ConfigEntity, webhook *EmbyWebhook) error {
	if webhook.Item.Type != "Episode" {
		return nil
	}

	// Check if the series name matches the webhook item series name
	itemEpisode, err := GetItem(ctx, cfg, webhook.Item.Id)
	if err != nil {
		return fmt.Errorf("failed to get Emby episode item: %w", err)
	}
	if itemEpisode.Type != "Episode" {
		return fmt.Errorf("item is not an episode: %s", itemEpisode.Name)
	}
	// Check if the series name matches the webhook item series name
	itemSeason, err := GetItem(ctx, cfg, itemEpisode.ParentId)
	if err != nil {
		return fmt.Errorf("failed to get Emby season item: %w", err)
	}
	if itemSeason.Type != "Season" {
		return fmt.Errorf("item is not a season: %s", itemSeason.Name)
	}
	// Check if the series name matches the webhook item series name
	itemSerie, err := GetItem(ctx, cfg, itemSeason.ParentId)
	if err != nil {
		return fmt.Errorf("failed to get Emby series item: %w", err)
	}
	if itemSerie.Type != "Series" {
		return fmt.Errorf("item is not a series: %s", itemSerie.Name)
	}
	imdbId, err := itemSerie.GetImdbId()
	if err != nil {
		return err
	}

	fmt.Printf("Processing Emby series: %s with IMDb ID: %s\n", webhook.Item.Name, imdbId)

	seasonNumber := int16(*webhook.Item.ParentIndexNumber)
	episodeNumber := int16(*webhook.Item.IndexNumber)

	traktRequest := trakt.MarkAsWatchedMap{}
	if err := traktRequest.AppendTvShow(imdbId, seasonNumber, episodeNumber, webhook.Date); err != nil {
		return fmt.Errorf("failed to append Emby series to trakt request: %w", err)
	}

	// Mark the item as watched in Trakt
	if err := trakt.MarkItemAsWatched(ctx, cfg, &traktRequest, false); err != nil {
		return err
	}
	fmt.Printf("Marked Emby series: %s as watched in Trakt with IMDB: %s\n", webhook.Item.Name, imdbId)

	return nil
}
