package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"trakt-sync/internal/config"
	"trakt-sync/internal/ctxutils"
	"trakt-sync/internal/emby"
	"trakt-sync/internal/trakt"
	"trakt-sync/internal/utils"

	"github.com/jarcoal/httpmock"

	_ "github.com/mattn/go-sqlite3"
)

func TestEmbyWebhook_MarkMovieAsWatchedOnTrakt(t *testing.T) {
	// Create a temporary SQLite database file
	tempDBFile := "./TestEmbyWebhook_MarkMovieAsWatchedOnTrakt.db"
	db, err := sql.Open("sqlite3", tempDBFile)
	if err != nil {
		t.Fatalf("Failed to create temporary database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(tempDBFile)
	}()

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxutils.ContextDbKey, db)
	config.InitConfigTable(&ctx)

	// Mock configuration
	cfg := config.ConfigEntity{
		Emby: &config.EmbyConfig{
			BaseURL: "http://localhost:8096",
			APIKey:  "test-api-key",
			UserID:  "aac3a78d9f184ea480fb1629e76aad57",
		},
		Trakt: &config.TraktConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			Code:         "test-code",
		},
	}
	err = config.UpsertConfig(&ctx, &cfg)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Read webhook data from file
	webhookBytes, err := os.ReadFile("./../../testdata/emby/webhooks/movies/mark_played.json")
	if err != nil {
		t.Fatalf("Failed to read webhook test file: %v", err)
	}

	// Unmarshal webhook data
	var webhook emby.EmbyWebhook
	if err := json.Unmarshal(webhookBytes, &webhook); err != nil {
		t.Fatalf("Failed to unmarshal webhook data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/emby/webhooks", bytes.NewReader(webhookBytes))
	req.Header.Set("Content-Type", "application/json")

	// Mock response recorder
	resp := httptest.NewRecorder()

	httpmock.Activate(t)
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", trakt.TraktApiUrl+"/sync/history",
		func(req *http.Request) (*http.Response, error) {

			request, err := utils.SerializeBody[trakt.MarkAsWatchedRequest](req.Body)
			if err != nil {
				t.Fatalf("Failed to serialize request body: %v", err)
			}
			if len(request.Movies) != 1 {
				t.Fatalf("Expected 1 movie in request, got %d", len(request.Movies))
			}
			if len(request.Shows) != 0 {
				t.Fatalf("Expected no shows in request, got %d", len(request.Shows))
			}
			imdb, err := webhook.GetImdbId()
			if err != nil || request.Movies[0].Ids.Imdb != imdb {
				t.Fatalf("Expected IMDB ID '%s', got '%s'", imdb, request.Movies[0].Ids.Imdb)
			}

			response := httpmock.NewStringResponse(http.StatusOK, `{"message": "Success"}`)
			return response, nil
		})

	// Call the handler
	HandleEmbyWebhooks(&ctx, resp, req)

	// Assert response
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}
}
