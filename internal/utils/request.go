package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"trakt-sync/internal/config"
)

// HttpGet performs a GET request and decodes the response into the generic type T.
// It accepts a function parameter to add headers to the request.
func HttpGet[T any](url string, config *config.ConfigEntity, addHeaders func(req *http.Request, config *config.ConfigEntity)) (*T, error) {
	// Construct the HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request to Emby: %w", err)
	}

	if addHeaders != nil {
		// Add headers to the request using the provided function
		addHeaders(req, config)
	}

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Decode the JSON response into the generic type
	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func HttpPost[TReq any, TRes any](url string, config *config.ConfigEntity, body *TReq, addHeaders func(req *http.Request, config *config.ConfigEntity)) (*TRes, error) {

	var reqBody *bytes.Buffer = nil

	// Check if body is nil before marshaling
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(b)
	}

	// Construct the HTTP request
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	if addHeaders != nil {
		// Add headers to the request using the provided function
		addHeaders(req, config)
	}

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Decode the JSON response into the generic type
	var result TRes
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
