package emby

import (
	"context"
	"fmt"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

func MarkItemAsWatched(ctx *context.Context, cfg *config.ConfigEntity, itemId string) error {

	if !cfg.Emby.IsValid(&config.EmbyOptions{}) {
		return fmt.Errorf("Emby configuration is invalid")
	}

	preUrl := "%s/Users/%s/PlayedItems/%s"
	url := fmt.Sprintf(preUrl, cfg.Emby.BaseURL, cfg.Emby.UserID, itemId)

	_, err := utils.HttpPost[struct{}, struct{}](
		utils.RequestParams{
			URL:     url,
			Config:  cfg,
			Context: ctx,
		},
		&struct{}{},
	)
	if err != nil {
		return fmt.Errorf("failed to mark emby item as watched: %w", err)
	}

	return nil
}
