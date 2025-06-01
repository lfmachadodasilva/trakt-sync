package trakt

import (
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

func Auth(c *config.ConfigEntity, code string) error {

	if code == "" {
		return fmt.Errorf("Trakt code is not set")
	}

	preUrl := "%s/oauth/token"
	url := fmt.Sprintf(preUrl, TraktApiUrl)

	authRequest := TraktAuthRequest{
		Code:         code,
		ClientID:     c.Trakt.ClientID,
		ClientSecret: c.Trakt.ClientSecret,
		RedirectURI:  c.Trakt.RedirectURL,
		GrantType:    "authorization_code",
	}

	res, err := utils.HttpPost[TraktAuthRequest, TraktAuthResponse](url, c, &authRequest, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch Emby items: %w", err)
	}

	err2 := config.UpsertConfig(&config.ConfigEntity{
		Trakt: &config.TraktConfig{
			ClientID:     c.Trakt.ClientID,
			ClientSecret: c.Trakt.ClientSecret,
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			Code:         code,
			RedirectURL:  c.Trakt.RedirectURL,
		},
	})
	if err2 != nil {
		return fmt.Errorf("failed to upsert config: %w", err2)
	}

	return nil
}

func AuthRefreshAccessToken(c *config.ConfigEntity) error {

	if c.Trakt == nil || c.Trakt.RefreshToken == "" {
		return fmt.Errorf("Trakt RefreshToken is not set")
	}
	preUrl := "%s/oauth/revoke"
	url := fmt.Sprintf(preUrl, TraktApiUrl)

	authRequest := TraktAuthRefreshRequest{
		Token:        c.Trakt.AccessToken,
		ClientID:     c.Trakt.ClientID,
		ClientSecret: c.Trakt.ClientSecret,
	}

	_, err := utils.HttpPost[TraktAuthRefreshRequest, struct{}](url, c, &authRequest, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch Emby items: %w", err)
	}

	return nil
}
