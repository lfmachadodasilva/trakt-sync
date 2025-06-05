package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	return SerializeBody[T](resp.Body)
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
	req, err := http.NewRequest(http.MethodPost, params.URL, reqBody)
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

	return SerializeBody[TRes](resp.Body)
}

// HttpError represents an HTTP error with a status code.
type HttpError struct {
	StatusCode int
	Err        error
}

// Error implements the error interface for HttpError.
func (e *HttpError) Error() string {
	return e.Err.Error()
}

// Unwrap provides compatibility with errors.Unwrap.
func (e *HttpError) Unwrap() error {
	return e.Err
}

// SerializeBody is a generic function that deserializes the body of an HTTP request into the specified type.
func SerializeBody[T any](body io.ReadCloser) (*T, error) {
	defer body.Close()

	var obj T
	if err := json.NewDecoder(body).Decode(&obj); err != nil {
		return nil, &HttpError{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	return &obj, nil
}
