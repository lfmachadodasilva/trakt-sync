package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
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

func TestEmbyMovieWebhook(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name            string
		webhookFile     string
		expectedStatus  int
		expectTraktCall bool
	}{
		{
			name:            "Paused_DontUpdateTrakt",
			webhookFile:     "./../../testdata/emby/webhooks/movies/pause.json",
			expectedStatus:  http.StatusOK,
			expectTraktCall: false,
		},
		{
			name:            "MarkPlayed_MarkMovieAsWatchedOnTrakt",
			webhookFile:     "./../../testdata/emby/webhooks/movies/mark_played.json",
			expectedStatus:  http.StatusOK,
			expectTraktCall: true,
		},
		{
			name:            "StopDone_MarkMovieAsWatchedOnTrakt",
			webhookFile:     "./../../testdata/emby/webhooks/movies/stop_done.json",
			expectedStatus:  http.StatusOK,
			expectTraktCall: true,
		},
		{
			name:            "StopNotDone_MarkMovieAsWatchedOnTrakt",
			webhookFile:     "./../../testdata/emby/webhooks/movies/stop_not_done.json",
			expectedStatus:  http.StatusOK,
			expectTraktCall: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()

			ctx := context.Background()

			// Create a temporary SQLite database file
			tempDBFile := "./TestEmbyWebhook_" + tc.name + ".db"
			db, err := sql.Open("sqlite3", tempDBFile)
			if err != nil {
				t.Fatalf("Failed to create temporary database: %v", err)
			}
			defer func() {
				db.Close()
				os.Remove(tempDBFile)
			}()
			// Initialize context and database
			testInitDatabase(t, &ctx, db)

			// Read webhook data from file and unmarshal it
			webhook, webhookBytes := unmarshalWebhook(t, tc.webhookFile)
			req := httptest.NewRequest(http.MethodPost, "/emby/webhooks", bytes.NewReader(webhookBytes))
			req.Header.Set("Content-Type", "application/json")

			// Mock response recorder
			resp := httptest.NewRecorder()

			// Activate httpmock for the current test
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Track if Trakt API call was made
			var traktCallMade bool = false

			// Register mock responder
			httpmock.RegisterResponder("POST", trakt.TraktApiUrl+"/sync/history",
				func(req *http.Request) (*http.Response, error) {

					traktCallMade = true

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

					jsonData, err := json.Marshal(trakt.MarkAsWatchedResponse{
						Added: struct {
							Movies   int16 `json:"movies,omitempty"`
							Episodes int16 `json:"episodes,omitempty"`
						}{
							Movies:   1,
							Episodes: 0,
						},
					})
					if err != nil {
						t.Fatalf("Failed to marshal configs: %v", err)
					}

					response := httpmock.NewStringResponse(http.StatusOK, `{"message": "Success"}`)
					response.Body = io.NopCloser(bytes.NewBuffer(jsonData))

					return response, nil
				})

			// Call the handler
			HandleEmbyWebhooks(&ctx, resp, req)

			// Assert response status
			if resp.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.Code)
			}
			// Assert Trakt API call
			if tc.expectTraktCall && !traktCallMade {
				t.Errorf("Expected Trakt API call, but it wasn't made")
			}
			if !tc.expectTraktCall && traktCallMade {
				t.Errorf("Expected no Trakt API call, but it was made")
			}
		})
	}
}

func TestEmbyShowWebhook(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name            string
		webhookFile     string
		expectedStatus  int
		expectTraktCall bool
	}{
		{
			name:            "MarkPlayed_MarkShowAsWatchedOnTrakt",
			webhookFile:     "./../../testdata/emby/webhooks/shows/item.markplayed.json",
			expectedStatus:  http.StatusOK,
			expectTraktCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()

			ctx := context.Background()

			// Create a temporary SQLite database file
			tempDBFile := "./TestEmbyWebhook_" + tc.name + ".db"
			db, err := sql.Open("sqlite3", tempDBFile)
			if err != nil {
				t.Fatalf("Failed to create temporary database: %v", err)
			}
			defer func() {
				db.Close()
				os.Remove(tempDBFile)
			}()
			cfg := testInitDatabase(t, &ctx, db)

			// Read webhook data from file and unmarshal it
			webhook, webhookBytes := unmarshalWebhook(t, tc.webhookFile)
			req := httptest.NewRequest(http.MethodPost, "/emby/webhooks", bytes.NewReader(webhookBytes))
			req.Header.Set("Content-Type", "application/json")

			// Mock response recorder
			resp := httptest.NewRecorder()

			// Activate httpmock for the current test
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Track if Trakt API call was made
			var traktCallMade bool = false

			var imdbSerie string = "tt7777" // Example IMDB ID for the series

			// Register mock responder
			httpmock.RegisterResponder("POST", trakt.TraktApiUrl+"/sync/history",
				func(req *http.Request) (*http.Response, error) {

					traktCallMade = true

					request, err := utils.SerializeBody[trakt.MarkAsWatchedRequest](req.Body)
					if err != nil {
						t.Fatalf("Failed to serialize request body: %v", err)
					}
					if len(request.Movies) != 0 {
						t.Fatalf("Expected no movie in request, got %d", len(request.Movies))
					}
					if len(request.Shows) != 1 {
						t.Fatalf("Expected 1 show in request, got %d", len(request.Shows))
					}
					if err != nil || request.Shows[0].Ids.Imdb != imdbSerie {
						t.Fatalf("Expected IMDB ID '%s', got '%s'", imdbSerie, request.Shows[0].Ids.Imdb)
					}

					if request.Shows[0].Seasons[0].Number != int16(*webhook.Item.ParentIndexNumber) {
						t.Fatalf("Expected season number %d, got %d", *webhook.Item.ParentIndexNumber, request.Shows[0].Seasons[0].Number)
					}
					if request.Shows[0].Seasons[0].Episodes[0].Number != int16(*webhook.Item.IndexNumber) {
						t.Fatalf("Expected episode number %d, got %d", *webhook.Item.IndexNumber, request.Shows[0].Seasons[0].Episodes[0].Number)
					}

					jsonData, err := json.Marshal(trakt.MarkAsWatchedResponse{
						Added: struct {
							Movies   int16 `json:"movies,omitempty"`
							Episodes int16 `json:"episodes,omitempty"`
						}{
							Movies:   0,
							Episodes: 1,
						},
					})
					if err != nil {
						t.Fatalf("Failed to marshal configs: %v", err)
					}

					response := httpmock.NewStringResponse(http.StatusOK, `{"message": "Success"}`)
					response.Body = io.NopCloser(bytes.NewBuffer(jsonData))

					return response, nil
				})

			var funcRegisterGetItem = func(itemId string, parentItemId string, itemType string, imdbId string) {
				httpmock.RegisterResponder("GET", cfg.Emby.BaseURL+"/Users/"+cfg.Emby.UserID+"/Items/"+itemId+"?Fields=ProviderIds",
					func(req *http.Request) (*http.Response, error) {
						response := httpmock.NewStringResponse(http.StatusOK, `{"message": "Success"}`)
						jsonData, err := json.Marshal(emby.EmbyItemResponse{
							ProviderIds: emby.EmbyProviderIds{
								Imdb: imdbId,
								IMDB: imdbId,
							},
							ParentId: parentItemId,
							Type:     itemType,
						})
						if err != nil {
							t.Fatalf("Failed to marshal configs: %v", err)
						}
						response.Body = io.NopCloser(bytes.NewBuffer(jsonData))
						return response, nil
					})
			}
			funcRegisterGetItem(webhook.Item.Id, "SeasonId", "Episode", imdbSerie)
			funcRegisterGetItem("SeasonId", "SerieId", "Season", imdbSerie)
			funcRegisterGetItem("SerieId", "Serie", "Series", imdbSerie)

			// Call the handler
			HandleEmbyWebhooks(&ctx, resp, req)

			// Assert response status
			if resp.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.Code)
			}
			// Assert Trakt API call
			if tc.expectTraktCall && !traktCallMade {
				t.Errorf("Expected Trakt API call, but it wasn't made")
			}
			if !tc.expectTraktCall && traktCallMade {
				t.Errorf("Expected no Trakt API call, but it was made")
			}
		})
	}
}

func unmarshalWebhook(t *testing.T, filePath string) (*emby.EmbyWebhook, []byte) {
	// Read webhook data from file
	webhookBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read webhook test file: %v", err)
	}
	// Unmarshal webhook data
	var webhook emby.EmbyWebhook
	if err := json.Unmarshal(webhookBytes, &webhook); err != nil {
		t.Fatalf("Failed to unmarshal webhook data: %v", err)
	}

	return &webhook, webhookBytes
}

func testInitDatabase(t *testing.T, ctx *context.Context, db *sql.DB) *config.ConfigEntity {
	// Initialize context and database
	*ctx = context.WithValue(*ctx, ctxutils.ContextDbKey, db)
	config.InitConfigTable(ctx)

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
	err := config.UpsertConfig(ctx, &cfg)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	return &cfg
}
