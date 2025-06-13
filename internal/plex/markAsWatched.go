package plex

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"trakt-sync/internal/ctxutils"

	"github.com/LukeHagar/plexgo"
)

func MarkAsWatched(ctx *context.Context, itemId string) error {
	plexSdk, ok := (*ctx).Value(ctxutils.ContextPlexSdkKey).(*plexgo.PlexAPI)
	if !ok {
		return fmt.Errorf("failed to get plex sdk from context")
	}

	itemIdFloat, err := strconv.ParseFloat(itemId, 64)
	if err != nil {
		log.Printf("Failed to convert itemId to float: %v", err)
		return err
	}

	res, err := plexSdk.Media.MarkPlayed(*ctx, itemIdFloat)
	if err != nil {
		log.Printf("Failed to mark item as played: %v", err)
		return err
	}

	if res != nil {
		log.Printf("Failed to mark item as watched: %v", err)
		return err
	}

	return nil
}
