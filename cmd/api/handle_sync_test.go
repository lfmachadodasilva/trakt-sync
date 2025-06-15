package main

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"trakt-sync/internal/config"
	"trakt-sync/internal/ctxutils"

	"github.com/jarcoal/httpmock"
)

func TestSyncEmby(t *testing.T) {

	// Create a temporary SQLite database file
	tempDBFile := "./TestSyncEmby.db"
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

	req := httptest.NewRequest(http.MethodPost, "/sync", nil)
	req.Header.Set("Content-Type", "application/json")

	// Mock response recorder
	resp := httptest.NewRecorder()

	// Activate httpmock for the current test
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	HandleSyncAll(&ctx, resp, req)

	// Assert response status
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestSyncPlex(t *testing.T) {

	// Create a temporary SQLite database file
	tempDBFile := "./TestSyncPlex.db"
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

	// Mock config database
	testInitDatabase(t, &ctx, db)
}
