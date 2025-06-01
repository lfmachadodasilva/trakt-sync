package models

type TraktConfig struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Code         string `json:"code,omitempty"`
	RedirectURL  string `json:"redirect_url,omitempty"`
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

type Config struct {
	Trakt    *TraktConfig    `json:"trakt,omitempty"`
	Emby     *EmbyConfig     `json:"emby,omitempty"`
	Plex     *PlexConfig     `json:"plex,omitempty"`
	Jellyfin *JellyfinConfig `json:"jellyfin,omitempty"`
}
