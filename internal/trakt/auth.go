package trakt

import (
	"context"
	"fmt"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

type TraktAuthRequest struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type TraktAuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type TraktAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	CreatedAt    int    `json:"created_at"`
	Scope        string `json:"scope"`
}

func Auth(ctx *context.Context, cfg *config.ConfigEntity, code string) error {

	if code == "" {
		return fmt.Errorf("Trakt code is not set")
	}

	preUrl := "%s/oauth/token"
	url := fmt.Sprintf(preUrl, TraktApiUrl)

	authRequest := TraktAuthRequest{
		Code:         code,
		ClientID:     cfg.Trakt.ClientID,
		ClientSecret: cfg.Trakt.ClientSecret,
		RedirectURI:  cfg.Trakt.RedirectURL,
		GrantType:    "authorization_code",
	}

	res, err := utils.HttpPost[TraktAuthRequest, TraktAuthResponse](
		utils.RequestParams{
			URL:    url,
			Config: cfg,
		},
		&authRequest,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch Emby items: %w", err)
	}

	err2 := config.UpsertConfig(ctx, &config.ConfigEntity{
		Trakt: &config.TraktConfig{
			ClientID:     cfg.Trakt.ClientID,
			ClientSecret: cfg.Trakt.ClientSecret,
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			Code:         code,
			RedirectURL:  cfg.Trakt.RedirectURL,
		},
	})
	if err2 != nil {
		return fmt.Errorf("failed to upsert config: %w", err2)
	}

	return nil
}

func AuthRefreshAccessToken(ctx *context.Context, cfg *config.ConfigEntity) error {

	if cfg.Trakt == nil || cfg.Trakt.RefreshToken == "" {
		return fmt.Errorf("Trakt RefreshToken is not set")
	}
	preUrl := "%s/oauth/token"
	url := fmt.Sprintf(preUrl, TraktApiUrl)

	authRequest := TraktAuthRefreshRequest{
		RefreshToken: cfg.Trakt.RefreshToken,
		ClientID:     cfg.Trakt.ClientID,
		ClientSecret: cfg.Trakt.ClientSecret,
		RedirectURI:  cfg.Trakt.RedirectURL,
		GrantType:    "refresh_token",
	}

	response, err := utils.HttpPost[TraktAuthRefreshRequest, TraktAuthResponse](
		utils.RequestParams{
			URL:    url,
			Config: cfg,
		},
		&authRequest,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch Emby items: %w", err)
	}

	// After revoking the token, we need to clear the access token and refresh token
	cfg.Trakt.AccessToken = response.AccessToken
	err2 := config.UpsertConfig(ctx, cfg)
	if err2 != nil {
		return fmt.Errorf("failed to upsert config: %w", err2)
	}

	return nil
}
