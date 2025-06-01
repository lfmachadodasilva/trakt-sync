package emby

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"trakt-sync/internal/models"
)

// EmbyUserResponse represents the structure of the user response from Emby
type EmbyUserResponse struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	ServerId string `json:"ServerId"`
	Prefix   string `json:"Prefix"`
}

// FetchEmbyUsers fetches user information from Emby using the provided models.Config
func FetchEmbyUsers(config *models.Config) ([]EmbyUserResponse, error) {
	// Get the Emby base URL from the config
	baseURL := config.Emby.BaseURL

	// Validate the Emby base URL
	if config.Emby == nil || config.Emby.BaseURL == "" {
		return nil, fmt.Errorf("Emby base URL is not configured")
	}

	// Check if the base URL is a valid URL
	_, err := url.ParseRequestURI(config.Emby.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("Emby base URL is invalid: %v", err)
	}

	// Construct the URL for the GET request
	url := fmt.Sprintf("%s/Users", baseURL)

	// Construct the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET /users request to Emby: %w", err)
	}

	// Add headers to the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Emby-Token", config.Emby.APIKey)

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /users request to Emby: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response from Emby: %d", resp.StatusCode)
	}

	// Decode the JSON response
	var users []EmbyUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode Emby response: %w", err)
	}

	return users, nil
}
