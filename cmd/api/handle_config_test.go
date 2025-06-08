package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"trakt-sync/internal/config"
	"trakt-sync/internal/ctxutils"
	"trakt-sync/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

func TestConfigGetEndpoint(t *testing.T) {

	// Create a temporary SQLite database file
	tempDBFile := "./TestConfigGetEndpoint.db"
	db, err := sql.Open("sqlite3", tempDBFile)
	if err != nil {
		t.Fatalf("Failed to create temporary database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(tempDBFile)
	}()

	// Initialize context and database
	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxutils.ContextDbKey, db)
	config.InitConfigTable(&ctx)

	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	req.Header.Set("Content-Type", "application/json")

	// Mock response recorder
	resp := httptest.NewRecorder()

	// Call the handler
	handleGetConfig(&ctx, resp)

	// Assert response status
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}

	restResponse := resp.Result()
	if restResponse == nil {
		t.Error("Expected non-nil response, got nil")
		return
	}

	// Assert response
	obj, err := utils.SerializeBody[config.ConfigEntity](restResponse.Body)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
		return
	}
	if obj.Emby.BaseURL == "" {
		t.Error("Expected 'emby' key in response body, got nil")
	}
	if obj.Trakt.RedirectURL == "" {
		t.Error("Expected 'trakt' key in response body, got nil")
	}
}

func TestConfigPatchEndpoint(t *testing.T) {

	// Create a temporary SQLite database file
	tempDBFile := "./TestConfigPatchEndpoint.db"
	db, err := sql.Open("sqlite3", tempDBFile)
	if err != nil {
		t.Fatalf("Failed to create temporary database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(tempDBFile)
	}()

	// Initialize context and database
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

	req := httptest.NewRequest(http.MethodPatch, "/config", nil)
	req.Header.Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println("Failed to marshal users:", err)
		return
	}
	req.Body = io.NopCloser(bytes.NewReader(jsonData))

	// Mock response recorder
	resp := httptest.NewRecorder()

	// Call the handler
	handlePatchConfig(&ctx, resp, req)

	// Assert response status
	if resp.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, resp.Code)
	}

	cfg2, err := config.ReadConfig(&ctx)
	if err != nil {
		t.Errorf("Failed to read config: %v", err)
		return
	}

	// Compare all fields between cfg and cfg2 using reflection
	compareConfigs := func(t *testing.T, name string, a, b interface{}) {
		va := reflect.ValueOf(a)
		vb := reflect.ValueOf(b)
		if va.IsNil() || vb.IsNil() {
			t.Errorf("Expected non-nil %s config", name)
			return
		}
		va = va.Elem()
		vb = vb.Elem()
		typ := va.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fa := va.Field(i)
			fb := vb.Field(i)
			if !reflect.DeepEqual(fa.Interface(), fb.Interface()) {
				t.Errorf("Expected %s.%s to be '%v', got '%v'", name, field.Name, fa.Interface(), fb.Interface())
			}
		}
	}

	compareConfigs(t, "Emby", cfg2.Emby, cfg.Emby)
	compareConfigs(t, "Trakt", cfg2.Trakt, cfg.Trakt)
	// compareConfigs(t, "Plex", cfg2.Plex, cfg.Plex)
	// compareConfigs(t, "Jellyfin", cfg2.Jellyfin, cfg.Jellyfin)
}
