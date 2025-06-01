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

type TraktAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func FetchTraktAuth(c *config.ConfigEntity, code string) error {

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
