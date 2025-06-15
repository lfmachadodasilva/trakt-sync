package trakt

import (
	"testing"
	"time"
)

func TestTraktAppendMovie(t *testing.T) {
	// Arrange
	traktMap := MarkAsWatchedMap{}
	time1 := time.Now().Add(-24 * time.Hour) // Set a time 24 hours ago
	time2 := time.Now()
	imdbID := "tt1234567"

	// Test adding a movie
	if err := traktMap.AppendMovie(imdbID, time1); err != nil {
		t.Errorf("Failed to append movie: %v", err)
	}
	if err := traktMap.AppendMovie(imdbID, time2); err != nil {
		t.Errorf("Failed to append movie: %v", err)
	}

	// Assert that the movie was added
	if time, exists := traktMap.Movies[imdbID]; !exists {
		t.Errorf("Expected movie '%s' to be added to traktMap, but it was not found", imdbID)
	} else if time != time1 {
		t.Errorf("Expected movie '%s' to have watchedAt time %v, but got %v", imdbID, time1, time)
	}
}

func TestTraktAppendTvShow(t *testing.T) {
	// Arrange
	traktMap := MarkAsWatchedMap{}
	imdbID := "tt1234567"
	seasonNumber := int16(1)
	episodeNumber := int16(1)
	time1 := time.Now().Add(-24 * time.Hour) // Set a time 24 hours ago
	time2 := time.Now()

	// Test adding a TV show episode
	if err := traktMap.AppendTvShow(imdbID, seasonNumber, episodeNumber, time1); err != nil {
		t.Errorf("Failed to append TV show: %v", err)
	}
	if err := traktMap.AppendTvShow(imdbID, seasonNumber, episodeNumber, time2); err != nil {
		t.Errorf("Failed to append TV show: %v", err)
	}

	// Assert that the TV show episode was added
	if time, exist := traktMap.Shows[imdbID][seasonNumber][episodeNumber]; !exist {
		t.Errorf("Expected episode %d of season %d for TV show '%s' to be added, but it was not found", episodeNumber, seasonNumber, imdbID)
	} else if time != time1 {
		t.Errorf("Expected episode %d of season %d for TV show '%s' to have watchedAt time %v, but got %v", episodeNumber, seasonNumber, imdbID, time1, time)
	}
}
