package plex

import (
	"context"
	"fmt"
	"trakt-sync/internal/config"
	"trakt-sync/internal/utils"
)

func MarkAsWatched(ctx *context.Context, cfg *config.ConfigEntity, itemId string) error {
	// plex sdk is not used and does not work to mark items as watched

	// plexSdk, ok := (*ctx).Value(ctxutils.ContextPlexSdkKey).(*plexgo.PlexAPI)
	// if !ok {
	// 	return fmt.Errorf("failed to get plex sdk from context")
	// }

	// itemIdFloat, err := strconv.ParseFloat(itemId, 64)
	// if err != nil {
	// 	log.Printf("Failed to convert itemId to float: %v", err)
	// 	return err
	// }

	// res, err := plexSdk.Media.MarkPlayed(*ctx, itemIdFloat)
	// if err != nil {
	// 	log.Printf("Failed to mark item as played: %v", err)
	// 	return err
	// }

	// if res != nil {
	// 	log.Printf("Failed to mark item as watched: %v", err)
	// 	return err
	// }

	preUrl := "%s/:/scrobble?key=%s&identifier=com.plexapp.plugins.library"
	url := fmt.Sprintf(preUrl, cfg.Plex.BaseURL, itemId)

	_, err := utils.HttpGet[struct{}](
		utils.RequestParams{
			URL:            url,
			Config:         cfg,
			AddHeaders:     addPlexHeaders,
			Context:        ctx,
			IgnoreResponse: true,
		})
	if err != nil {
		return fmt.Errorf("failed to mark plex item as watched: %w", err)
	}

	return nil
}
