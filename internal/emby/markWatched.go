package emby

import (
	"fmt"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

func MarkItemAsWatched(c *config.ConfigEntity, itemId string) error {

	if !c.Emby.IsValid(&config.EmbyOptions{}) {
		return fmt.Errorf("Emby configuration is invalid")
	}

	preUrl := "%s/Users/%s/PlayedItems/%s"
	url := fmt.Sprintf(preUrl, c.Emby.BaseURL, c.Emby.UserID, itemId)

	_, err := utils.HttpPost[struct{}, struct{}](
		utils.RequestParams{
			URL:    url,
			Config: c,
		},
		&struct{}{},
	)
	if err != nil {
		return fmt.Errorf("failed to mark item as watched: %w", err)
	}

	return nil
}
