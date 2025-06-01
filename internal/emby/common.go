package emby

import (
	"net/http"
	"trakt-sync/internal/config"
)

func addEmbyHeaders(req *http.Request, config *config.ConfigEntity) {
	req.Header.Set("X-Emby-Token", config.Emby.APIKey)
	req.Header.Set("Content-Type", "application/json")
}
