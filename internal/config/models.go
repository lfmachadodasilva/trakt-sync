package config

import (
	"net/url"
	"time"
)

type TraktConfig struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Code         string `json:"code,omitempty"`
	RedirectURL  string `json:"redirect_url,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	CreatedAt    int    `json:"created_at,omitempty"`
}

type EmbyConfig struct {
	UserID  string `json:"user_id,omitempty"`
	APIKey  string `json:"api_key,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
}

type PlexConfig struct {
	UserID string `json:"user_id,omitempty"`
}

type JellyfinConfig struct {
	UserID string `json:"user_id,omitempty"`
}

type ConfigEntity struct {
	Trakt    *TraktConfig    `json:"trakt,omitempty"`
	Emby     *EmbyConfig     `json:"emby,omitempty"`
	Plex     *PlexConfig     `json:"plex,omitempty"`
	Jellyfin *JellyfinConfig `json:"jellyfin,omitempty"`
}

func (emby *EmbyConfig) IsValid(ignoreUserId bool) bool {
	if emby == nil {
		return false
	}

	if ignoreUserId && emby.UserID == "" {
		return false
	}

	// Validate the Emby base URL
	if emby.BaseURL == "" || emby.APIKey == "" {
		return false
	}

	// Check if the base URL is a valid URL
	_, err := url.ParseRequestURI(emby.BaseURL)
	if err != nil {
		return false
	}

	return true
}

// IsAccessTokenValid checks if the access token is still valid
func (t *TraktConfig) IsAccessTokenValid() bool {
	if t.AccessToken == "" || t.CreatedAt == 0 || t.ExpiresIn == 0 {
		return false
	}

	// Calculate the expiration time
	expirationTime := time.Unix(int64(t.CreatedAt), 0).Add(time.Duration(t.ExpiresIn) * time.Second)

	// Compare with the current time
	return time.Now().Before(expirationTime)
}
