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
	APIKey  string `json:"api_key,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
}

type JellyfinConfig struct {
	APIKey  string `json:"api_key,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
}

type ConfigEntity struct {
	Trakt    *TraktConfig    `json:"trakt,omitempty"`
	Emby     *EmbyConfig     `json:"emby,omitempty"`
	Plex     *PlexConfig     `json:"plex,omitempty"`
	Jellyfin *JellyfinConfig `json:"jellyfin,omitempty"`
	Cronjob  string          `json:"cronjob,omitempty"`
}

type EmbyOptions struct {
	IgnoreUserId bool
}

type TraktOptions struct {
	IgnoreCode         bool
	IgnoreClientSecret bool
	IgnoreAccessToken  bool
}

func (emby *EmbyConfig) IsValid(options *EmbyOptions) bool {
	if emby == nil ||
		(options.IgnoreUserId && emby.UserID == "") ||
		emby.BaseURL == "" ||
		emby.APIKey == "" {
		return false
	}

	// Check if the base URL is a valid URL
	_, err := url.ParseRequestURI(emby.BaseURL)
	if err != nil {
		return false
	}

	return true
}

func (trakt *TraktConfig) IsValid(options *TraktOptions) bool {
	if trakt == nil ||
		(options.IgnoreCode && trakt.Code == "") ||
		(options.IgnoreClientSecret && trakt.ClientSecret == "") ||
		(options.IgnoreAccessToken && trakt.AccessToken == "") ||
		trakt.RedirectURL == "" ||
		trakt.ClientID == "" {
		return false
	}

	// Check if the base URL is a valid URL
	_, err := url.ParseRequestURI(trakt.RedirectURL)
	if err != nil {
		return false
	}

	return true
}

// IsAccessTokenValid checks if the access token is still valid
func (t *TraktConfig) IsAccessTokenValid() bool {
	if t.CreatedAt == 0 || t.ExpiresIn == 0 {
		return false
	}

	// Calculate the expiration time
	expirationTime := time.Unix(int64(t.CreatedAt), 0).Add(time.Duration(t.ExpiresIn) * time.Second)

	// Compare with the current time
	return time.Now().Before(expirationTime)
}
