package main

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"trakt-sync/internal/ctxutils"

	_ "github.com/mattn/go-sqlite3"
)

func TestConfigGetEndpoint(t *testing.T) {
	// Skip the test
	t.Skip("Skipping TestConfigGetEndpoint")

	// Create a temporary SQLite database file
	tempDBFile := "test.db"
	db, err := sql.Open("sqlite3", tempDBFile)
	if err != nil {
		t.Fatalf("Failed to create temporary database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(tempDBFile)
	}()

	// Initialize the database schema
	_, err = db.Exec(`CREATE TABLE config (key TEXT, value TEXT);`)
	if err != nil {
		t.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Seed the database with test data
	_, err = db.Exec(`INSERT INTO config (key, value) VALUES ('testKey', 'testValue');`)
	if err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}

	// Create a context with the test database
	ctx := context.WithValue(context.Background(), ctxutils.ContextDbKey, db)

	// Create a test HTTP server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleConfig(&ctx)(w, r)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Send a GET request to the /config endpoint
	resp, err := http.Get(ts.URL + "/config")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Validate the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Additional validation of the response body can be added here
}
