package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"trakt-sync/internal/config"
	"trakt-sync/internal/emby"
	"trakt-sync/internal/plex"
	"trakt-sync/internal/trakt"
)

var (
	spacePrefix1 = "*"
	spacePrefix2 = "**"
	spacePrefix3 = "***"
	spacePrefix4 = "****"
)

func HandleSync(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the sub-path after /sync/
		subPath := r.URL.Path[len("/sync"):]

		switch subPath {
		case "":
			// Handle the base /sync endpoint
			switch r.Method {
			case http.MethodPost:
				HandleSyncAll(ctx, w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Handle other sub-paths under /sync/
			http.Error(w, "Endpoint not found", http.StatusNotFound)
		}
	}
}

type SyncRequest struct {
	TraktItems      *trakt.TraktWatched
	TraktWatchedMap *trakt.MarkAsWatchedMap
	EmbyItems       *[]emby.EmbyItemResponse
	PlexItems       *plex.PlexItemsResponse
}

func HandleSyncAll(ctx *context.Context, w http.ResponseWriter, r *http.Request) {

	// read the config
	cfg, err := config.ReadConfig(ctx)
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusBadRequest)
		fmt.Println("Failed to read configs:", err)
		return
	}

	traktData, err := trakt.GetWatched(ctx, cfg)
	if err != nil {
		http.Error(w, "Failed to fetch watched data from Trakt: "+err.Error(), http.StatusBadRequest)
		log.Println("Failed to fetch watched data from Trakt:", err)
		return
	}

	traktWatchedMap := trakt.MarkAsWatchedMap{}

	if err := syncEmby(ctx, cfg, &traktData, &traktWatchedMap); err != nil {
		http.Error(w, "Failed to sync Emby with Trakt data: "+err.Error(), http.StatusBadRequest)
		fmt.Println("Failed to sync Emby with Trakt data:", err)
		return
	}

	if err := syncPlex(ctx, cfg, &traktData, &traktWatchedMap); err != nil {
		http.Error(w, "Failed to sync Plex with Trakt data: "+err.Error(), http.StatusBadRequest)
		log.Println("Failed to sync Plex with Trakt data:", err)
		return
	}

	if err := trakt.MarkItemAsWatched(ctx, cfg, &traktWatchedMap, false); err != nil {
		http.Error(w, "failed to mark items as watched in Trakt: "+err.Error(), http.StatusBadRequest)
		log.Println("Failed to sync Plex with Trakt data:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func syncEmby(ctx *context.Context, cfg *config.ConfigEntity, traktData *trakt.TraktWatched, traktWatchedMap *trakt.MarkAsWatchedMap) error {

	embyData, err := emby.GetAllItems(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to fetch Emby data: %w", err)
	}

	if err := syncEmbyMovie(ctx, cfg, &SyncRequest{
		TraktItems:      traktData,
		TraktWatchedMap: traktWatchedMap,
		EmbyItems:       &embyData.Movies,
	}); err != nil {
		return fmt.Errorf("failed to sync Emby movies: %w", err)
	}

	if err := syncEmbyTvShows(ctx, cfg, &SyncRequest{
		TraktItems:      traktData,
		TraktWatchedMap: traktWatchedMap,
		EmbyItems:       &embyData.Series,
	}); err != nil {
		return fmt.Errorf("failed to sync Emby shows: %w", err)
	}

	return nil
}

func syncEmbyMovie(ctx *context.Context, cfg *config.ConfigEntity, request *SyncRequest) error {

	fmt.Println("== Processing Emby items for movies ==")
	for _, embyItem := range *request.EmbyItems {
		imdbId := embyItem.ProviderIds.Imdb
		if imdbId == "" {
			imdbId = embyItem.ProviderIds.IMDB
			if imdbId == "" {
				fmt.Println("Skipping Emby item with no IMDb ID:", embyItem.Name)
				continue
			}
		}

		_, traktExists := request.TraktItems.Movies[imdbId]

		fmt.Println(spacePrefix1, "Processing Emby movie:", embyItem.Name, "with IMDb ID:", imdbId)
		if !traktExists && embyItem.UserData.Played {

			fmt.Println(spacePrefix2, "Adding new movie to Trakt request")

			if err := request.TraktWatchedMap.AppendMovie(imdbId, embyItem.UserData.LastPlayedDate); err != nil {
				return fmt.Errorf("failed to append Emby movie to Trakt request: %w", err)
			}

		} else if traktExists && !embyItem.UserData.Played {

			fmt.Println(spacePrefix2, "Marking Emby movie as watched")
			if err := emby.MarkItemAsWatched(ctx, cfg, embyItem.Id); err != nil {
				return fmt.Errorf("failed to mark Emby item as watched: %w", err)
			}

		} else {
			// Both are in sync, do nothing
			continue
		}
	}

	return nil
}

func syncEmbyTvShows(ctx *context.Context, cfg *config.ConfigEntity, request *SyncRequest) error {

	fmt.Println("== Processing Emby items for TV shows ==")
	for _, embyShow := range *request.EmbyItems {
		imdbId := embyShow.ProviderIds.Imdb
		if imdbId == "" {
			imdbId = embyShow.ProviderIds.IMDB
			if imdbId == "" {
				fmt.Println("Skipping Emby item with no IMDb ID:", embyShow.Name)
				continue
			}
		}

		fmt.Println(spacePrefix1, "Processing Emby show:", embyShow.Name, "with IMDb ID:", imdbId)
		for _, embyEpisode := range embyShow.Episodes {

			embySeasonNumber := embyEpisode.ParentIndexNumber
			embyEpisodeNumber := embyEpisode.IndexNumber

			if embySeasonNumber <= 0 || embyEpisodeNumber <= 0 {
				fmt.Println(spacePrefix2, "Skipping Emby episode with invalid season or episode number:", embyShow.Name, embySeasonNumber, embyEpisodeNumber)
				continue
			}

			_, traktImdbEpisodeExists := request.TraktItems.Shows[imdbId][embySeasonNumber][embyEpisodeNumber]
			if !traktImdbEpisodeExists && embyShow.UserData.Played {

				fmt.Printf("%s Marking Trakt episode as watched for S%dE%d\n", spacePrefix2, embySeasonNumber, embyEpisodeNumber)

				if err := request.TraktWatchedMap.AppendTvShow(imdbId, embySeasonNumber, embyEpisodeNumber, embyShow.UserData.LastPlayedDate); err != nil {
					return fmt.Errorf("failed to append Emby TV show to Trakt request: %w", err)
				}

			} else if traktImdbEpisodeExists && !embyShow.UserData.Played {

				fmt.Printf("%s Marking Emby episode as watched for S%dE%d\n", spacePrefix2, embySeasonNumber, embyEpisodeNumber)
				emby.MarkItemAsWatched(ctx, cfg, embyEpisode.Id)

			} else {
				// Both are in sync, do nothing
				continue
			}
		}
	}

	return nil
}

func syncPlex(ctx *context.Context, cfg *config.ConfigEntity, traktData *trakt.TraktWatched, traktWatchedMap *trakt.MarkAsWatchedMap) error {

	// Initialize plex & add client to context
	ctxCopy, err := plex.InitPlex(ctx, cfg)
	if err != nil {
		fmt.Println("Failed to initialize Plex:", err)
	}
	ctx = &ctxCopy

	plexItems, err := plex.GetAllItems(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to fetch plex items: %w", err)
	}

	if err := syncPlexMovies(ctx, cfg, &SyncRequest{
		TraktItems:      traktData,
		TraktWatchedMap: traktWatchedMap,
		PlexItems:       plexItems,
	}); err != nil {
		return fmt.Errorf("failed to sync Plex movies")
	}

	if err := syncPlexTvShows(ctx, cfg, &SyncRequest{
		TraktItems:      traktData,
		TraktWatchedMap: traktWatchedMap,
		PlexItems:       plexItems,
	}); err != nil {
		return fmt.Errorf("failed to sync Plex TV shows")
	}

	return nil
}

func syncPlexMovies(ctx *context.Context, cfg *config.ConfigEntity, request *SyncRequest) error {

	if request.PlexItems == nil || request.PlexItems.Movies == nil || len(request.PlexItems.Movies) == 0 {
		fmt.Println("No Plex movies to process")
		return nil
	}

	fmt.Println("== Processing Plex items for movies ==")
	for _, plexMovie := range request.PlexItems.Movies {

		_, traktExists := request.TraktItems.Movies[plexMovie.ImdbId]

		fmt.Println(spacePrefix1, "Processing Plex movie:", plexMovie.Title, "with IMDb ID:", plexMovie.ImdbId)
		if !traktExists && plexMovie.Watched {

			fmt.Println(spacePrefix2, "Adding new movie to Trakt request")
			if err := request.TraktWatchedMap.AppendMovie(plexMovie.ImdbId, plexMovie.WatchedTime); err != nil {
				return fmt.Errorf("failed to append Plex movie to Trakt request: %w", err)
			}

		} else if traktExists && !plexMovie.Watched {

			fmt.Println(spacePrefix2, "Marking Plex movie as watched")
			if err := plex.MarkAsWatched(ctx, cfg, plexMovie.Id); err != nil {
				fmt.Println(spacePrefix2, "Failed to mark Plex movie as watched:", err)
			}

		} else {
			// Both are in sync, do nothing
			continue
		}
	}

	return nil
}

func syncPlexTvShows(ctx *context.Context, cfg *config.ConfigEntity, request *SyncRequest) error {
	if request.PlexItems == nil || request.PlexItems.TvShows == nil || len(request.PlexItems.TvShows) == 0 {
		fmt.Println("No Plex TV shows to process")
		return nil
	}

	fmt.Println("== Processing Plex items for tv shows ==")
	for _, plexTvShow := range request.PlexItems.TvShows {

		for _, plexEpisode := range plexTvShow.Episodes {
			_, traktExists := request.TraktItems.Shows[plexTvShow.ImdbId][plexEpisode.Season][plexEpisode.Episode]

			fmt.Println(spacePrefix1, "Processing Plex episode:", plexTvShow.Title, "with IMDb ID:", plexTvShow.ImdbId)
			if !traktExists && plexEpisode.Watched {

				fmt.Println(spacePrefix2, "Adding new movie to Trakt request")
				if err := request.TraktWatchedMap.AppendTvShow(plexTvShow.ImdbId, plexEpisode.Season, plexEpisode.Episode, plexEpisode.WatchedTime); err != nil {
					return fmt.Errorf("failed to append Plex TV show to Trakt request: %w", err)
				}

			} else if traktExists && !plexEpisode.Watched {
				fmt.Println(spacePrefix2, "Marking Plex episode as watched")
				if err := plex.MarkAsWatched(ctx, cfg, plexEpisode.Id); err != nil {
					fmt.Println(spacePrefix2, "Failed to mark Plex episode as watched:", err)
				}
			} else {
				// Both are in sync, do nothing
				continue
			}
		}
	}

	return nil
}
