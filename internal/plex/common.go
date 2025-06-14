package plex

import (
	"net/http"
	"trakt-sync/internal/config"
)

func addPlexHeaders(req *http.Request, config *config.ConfigEntity) {
	req.Header.Set("X-Plex-Token", config.Plex.APIKey)
}
