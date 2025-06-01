package models

type TraktConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Code         string `json:"code"`
	RedirectURL  string `json:"redirect_url"`
}

type EmbyConfig struct {
	UserID  string `json:"user_id"`
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

type PlexConfig struct {
	UserID string `json:"user_id"`
}

type JellyfinConfig struct {
	UserID string `json:"user_id"`
}

type Config struct {
	Trakt    *TraktConfig    `json:"trakt,omitempty"`
	Emby     *EmbyConfig     `json:"emby,omitempty"`
	Plex     *PlexConfig     `json:"plex,omitempty"`
	Jellyfin *JellyfinConfig `json:"jellyfin,omitempty"`
}
