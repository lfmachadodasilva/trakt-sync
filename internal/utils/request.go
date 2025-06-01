package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trakt-sync/internal/models"
)

// Get performs a GET request and decodes the response into the generic type T.
// It accepts a function parameter to add headers to the request.
func Get[T any](url string, config *models.Config, addHeaders func(req *http.Request, config *models.Config)) (*T, error) {
	// Construct the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request to Emby: %w", err)
	}

	// Add headers to the request using the provided function
	addHeaders(req, config)

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to Emby: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response from Emby: %d", resp.StatusCode)
	}

	// Decode the JSON response into the generic type
	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Emby response: %w", err)
	}

	return &result, nil
}
