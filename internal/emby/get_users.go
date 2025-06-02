package emby

import (
	"fmt"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

// EmbyUserResponse represents the structure of the user response from Emby
type EmbyUserResponse struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	ServerId string `json:"ServerId"`
	Prefix   string `json:"Prefix"`
}

// FetchEmbyUsers fetches user information from Emby using the provided config.Config
func FetchEmbyUsers(cfg *config.ConfigEntity) ([]EmbyUserResponse, error) {
	// Validate the Emby configuration
	if !cfg.Emby.IsValid(&config.EmbyOptions{}) {
		return nil, fmt.Errorf("Emby configuration is invalid")
	}

	// Construct the URL for the GET request
	url := fmt.Sprintf("%s/Users", cfg.Emby.BaseURL)

	users, err := utils.HttpGet[[]EmbyUserResponse](
		utils.RequestParams{
			URL:        url,
			Config:     cfg,
			AddHeaders: addEmbyHeaders,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Emby users: %w", err)
	}

	return *users, nil
}
