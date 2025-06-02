package emby

import (
	"context"
	"fmt"
	"sync"
	"time"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

type EmbyItemsResponse struct {
	Items []EmbyItemResponse `json:"Items"`
}

type EmbyItemResponse struct {
	Id                string             `json:"Id"`
	Name              string             `json:"Name"`
	ServerId          string             `json:"ServerId"`
	Type              string             `json:"Type"`
	UserData          EmbyUserData       `json:"UserData"`
	ProviderIds       EmbyProviderIds    `json:"ProviderIds"`
	IndexNumber       int16              `json:"IndexNumber"`
	ParentIndexNumber int16              `json:"ParentIndexNumber"`
	Episodes          []EmbyItemResponse `json:"Episodes,omitempty"`
}

type EmbyUserData struct {
	Played         bool      `json:"Played"`
	LastPlayedDate time.Time `json:"LastPlayedDate"`
}

type EmbyProviderIds struct {
	Imdb string `json:"Imdb"`
	IMDB string `json:"IMDB"`
}

type EmbyItems struct {
	Movies []EmbyItemResponse
	Series []EmbyItemResponse
}

func GetAllItems(ctx *context.Context, cfg *config.ConfigEntity) (EmbyItems, error) {
	var movies []EmbyItemResponse
	var series []EmbyItemResponse
	var moviesErr, seriesErr error

	// Use WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(2) // Add two tasks to the WaitGroup

	// Fetch movies in a separate goroutine
	go func() {
		defer wg.Done() // Mark this task as done
		movies, moviesErr = GetItemsByType(ctx, cfg, "Movie")
	}()

	// Fetch series in a separate goroutine
	go func() {
		defer wg.Done() // Mark this task as done
		series, seriesErr = GetItemsByType(ctx, cfg, "Series")
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Check for errors
	if moviesErr != nil {
		return EmbyItems{}, fmt.Errorf("failed to fetch Emby movies: %w", moviesErr)
	}
	if seriesErr != nil {
		return EmbyItems{}, fmt.Errorf("failed to fetch Emby series: %w", seriesErr)
	}

	return EmbyItems{
		Movies: movies,
		Series: series,
	}, nil
}

func GetItemsByType(ctx *context.Context, cfg *config.ConfigEntity, itemType string) ([]EmbyItemResponse, error) {

	// Validate the itemType parameter
	if itemType != "Movie" && itemType != "Series" {
		return nil, fmt.Errorf("invalid itemType: %s. Must be 'Movie' or 'Series'", itemType)
	}

	// Validate the Emby configuration
	if !cfg.Emby.IsValid(&config.EmbyOptions{IgnoreUserId: true}) {
		return nil, fmt.Errorf("Emby configuration is invalid")
	}

	// Construct the URL for the GET request
	preUrl := "%s/Users/%s/Items?IncludeItemTypes=%s&Recursive=true&Fields=ProviderIds"
	url := fmt.Sprintf(preUrl, cfg.Emby.BaseURL, cfg.Emby.UserID, itemType)

	items, err := utils.HttpGet[EmbyItemsResponse](
		utils.RequestParams{
			URL:        url,
			Config:     cfg,
			AddHeaders: addEmbyHeaders,
			Context:    ctx,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Emby items: %w", err)
	}

	// If the itemType is "Series", fetch episodes for each series item
	if itemType == "Series" {
		for i := range items.Items {
			// Fetch episodes for each series item
			if items.Items[i].Id != "" {
				episodes, err := getEpisodes(ctx, cfg, &items.Items[i].Id)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch episodes for series %s: %w", items.Items[i].Name, err)
				}
				items.Items[i].Episodes = episodes
			}

			// Add cancellation check in every for loop
			select {
			case <-(*ctx).Done():
				return nil, (*ctx).Err() // Exit the loop if the context is canceled
			default:
				// Continue processing
			}
		}
	}

	return items.Items, nil
}

func getEpisodes(ctx *context.Context, cfg *config.ConfigEntity, embyId *string) ([]EmbyItemResponse, error) {

	// Validate the Emby configuration
	if !cfg.Emby.IsValid(&config.EmbyOptions{IgnoreUserId: true}) {
		return nil, fmt.Errorf("Emby configuration is invalid")
	}

	// Construct the URL for the GET request
	preUrl := "%s/Shows/%s/Episodes?&Recursive=true&EnableUserData=true&Fields=ProviderIds&UserId=%s"
	url := fmt.Sprintf(preUrl, cfg.Emby.BaseURL, *embyId, cfg.Emby.UserID)

	items, err := utils.HttpGet[EmbyItemsResponse](
		utils.RequestParams{
			URL:        url,
			Config:     cfg,
			AddHeaders: addEmbyHeaders,
			Context:    ctx,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Emby items: %w", err)
	}

	return items.Items, nil
}
