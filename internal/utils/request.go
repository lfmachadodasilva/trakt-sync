package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"trakt-sync/internal/config"
)

// RequestParams represents the parameters for HTTP requests
type RequestParams struct {
	Context    *context.Context
	URL        string
	Config     *config.ConfigEntity
	AddHeaders func(req *http.Request, config *config.ConfigEntity)
}

// HttpGet performs a GET request and decodes the response into the generic type T.
// It accepts a function parameter to add headers to the request.
func HttpGet[T any](params RequestParams) (*T, error) {
	// Construct the HTTP request
	req, err := http.NewRequestWithContext(*params.Context, http.MethodGet, params.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request to Emby: %w", err)
	}

	if params.AddHeaders != nil {
		// Add headers to the request using the provided function
		params.AddHeaders(req, params.Config)
	} else {
		req.Header.Set("Content-Type", "application/json")
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

func HttpPost[TReq any, TRes any](params RequestParams, body *TReq) (*TRes, error) {

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
	req, err := http.NewRequestWithContext(*params.Context, http.MethodPost, params.URL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	if params.AddHeaders != nil {
		// Add headers to the request using the provided function
		params.AddHeaders(req, params.Config)
	} else {
		req.Header.Set("Content-Type", "application/json")
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
