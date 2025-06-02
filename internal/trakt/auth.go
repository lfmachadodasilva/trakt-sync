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
	Token        string `json:"token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TraktAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	CreatedAt    int    `json:"created_at"`
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
	preUrl := "%s/oauth/revoke"
	url := fmt.Sprintf(preUrl, TraktApiUrl)

	authRequest := TraktAuthRefreshRequest{
		Token:        cfg.Trakt.AccessToken,
		ClientID:     cfg.Trakt.ClientID,
		ClientSecret: cfg.Trakt.ClientSecret,
	}

	_, err := utils.HttpPost[TraktAuthRefreshRequest, struct{}](
		utils.RequestParams{
			URL:    url,
			Config: cfg,
		},
		&authRequest,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch Emby items: %w", err)
	}

	return nil
}
