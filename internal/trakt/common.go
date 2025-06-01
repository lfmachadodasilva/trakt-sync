package trakt

import (
	"net/http"
	"trakt-sync/internal/config"
)

const (
	TraktApiUrl      = "https://api.trakt.tv"
	TraktApiVersion  = "2"
	TraktRedirectUri = "urn:ietf:wg:oauth:2.0:oob"
)

func addTraktHeaders(req *http.Request, config *config.ConfigEntity) {
	req.Header.Set("trakt-api-version", TraktApiVersion)
	req.Header.Set("trakt-api-key", config.Trakt.ClientID)
	req.Header.Set("Authorization", "Bearer "+config.Trakt.AccessToken)
	req.Header.Set("Content-Type", "application/json")
}
