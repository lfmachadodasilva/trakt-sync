package emby

import (
	"fmt"
	"sync"
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
	IndexNumber       int                `json:"IndexNumber"`
	ParentIndexNumber int                `json:"ParentIndexNumber"`
	Episodes          []EmbyItemResponse `json:"Episodes,omitempty"`
}

type EmbyUserData struct {
	Played bool `json:"Played"`
}

type EmbyProviderIds struct {
	Imdb string `json:"Imdb"`
	IMDB string `json:"IMDB"`
}

type EmbyItems struct {
	Movies []EmbyItemResponse
	Series []EmbyItemResponse
}

func GetAllItems(config *config.ConfigEntity) (EmbyItems, error) {
	var movies []EmbyItemResponse
	var series []EmbyItemResponse
	var moviesErr, seriesErr error

	// Use WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(2) // Add two tasks to the WaitGroup

	// Fetch movies in a separate goroutine
	go func() {
		defer wg.Done() // Mark this task as done
		movies, moviesErr = GetItemsByType(config, "Movie")
	}()

	// Fetch series in a separate goroutine
	go func() {
		defer wg.Done() // Mark this task as done
		series, seriesErr = GetItemsByType(config, "Series")
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

func GetItemsByType(config *config.ConfigEntity, itemType string) ([]EmbyItemResponse, error) {

	// Validate the itemType parameter
	if itemType != "Movie" && itemType != "Series" {
		return nil, fmt.Errorf("invalid itemType: %s. Must be 'Movie' or 'Series'", itemType)
	}

	// Validate the Emby configuration
	if !config.Emby.IsValid(true) {
		return nil, fmt.Errorf("Emby configuration is invalid")
	}

	// Construct the URL for the GET request
	preUrl := "%s/Users/%s/Items?IncludeItemTypes=%s&Recursive=true&Fields=ProviderIds"
	url := fmt.Sprintf(preUrl, config.Emby.BaseURL, config.Emby.UserID, itemType)

	items, err := utils.HttpGet[EmbyItemsResponse](url, config, addEmbyHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Emby items: %w", err)
	}

	// If the itemType is "Series", fetch episodes for each series item
	if itemType == "Series" {
		for i := range items.Items {
			// Fetch episodes for each series item
			if items.Items[i].Id != "" {
				episodes, err := getEpisodes(config, &items.Items[i].Id)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch episodes for series %s: %w", items.Items[i].Name, err)
				}
				items.Items[i].Episodes = episodes
			}
		}
	}

	return items.Items, nil
}

func getEpisodes(config *config.ConfigEntity, embyId *string) ([]EmbyItemResponse, error) {

	// Validate the Emby configuration
	if !config.Emby.IsValid(true) {
		return nil, fmt.Errorf("Emby configuration is invalid")
	}

	// Construct the URL for the GET request
	preUrl := "%s/Shows/%s/Episodes?&Recursive=true&EnableUserData=true&Fields=ProviderIds&UserId=%s"
	url := fmt.Sprintf(preUrl, config.Emby.BaseURL, *embyId, config.Emby.UserID)

	items, err := utils.HttpGet[EmbyItemsResponse](url, config, addEmbyHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Emby items: %w", err)
	}

	return items.Items, nil
}
