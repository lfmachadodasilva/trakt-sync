package main

import (
	"context"
	"fmt"
	"net/http"
	"trakt-sync/internal/config"
	"trakt-sync/internal/emby"
	"trakt-sync/internal/trakt"
	"trakt-sync/internal/utils"
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

func HandleSyncAll(ctx *context.Context, w http.ResponseWriter, r *http.Request) {

	cfg, err := config.ReadConfig(ctx)
	if err != nil {
		http.Error(w, "Failed to read configs", http.StatusBadRequest)
		return
	}

	traktData, err := trakt.GetWatched(ctx, &cfg)
	if err != nil {
		http.Error(w, "Failed to fetch watched data from Trakt: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := syncEmby(ctx, &cfg, traktData); err != nil {
		http.Error(w, "Failed to sync Emby with Trakt data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func syncEmby(ctx *context.Context, cfg *config.ConfigEntity, traktData *trakt.TraktWatched) error {

	embyData, err := emby.GetAllItems(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to fetch Emby data: %w", err)
	}

	traktRequest := &trakt.MarkAsWatchedRequest{}

	if err := syncEmbyMovie(ctx, cfg, &SyncRequest{
		TraktItems:   &traktData.Movies,
		TraktRequest: traktRequest,
		EmbyItems:    &embyData.Movies,
	}); err != nil {
		return fmt.Errorf("failed to sync Emby movies: %w", err)
	}

	if err := syncEmbyTvShows(ctx, cfg, &SyncRequest{
		TraktItems:   &traktData.Shows,
		TraktRequest: traktRequest,
		EmbyItems:    &embyData.Series,
	}); err != nil {
		return fmt.Errorf("failed to sync Emby shows: %w", err)
	}

	if err := trakt.MarkItemAsWatched(ctx, cfg, traktRequest); err != nil {
		return fmt.Errorf("failed to mark items as watched in Trakt: %w", err)
	}

	return nil
}

type SyncRequest struct {
	TraktItems   *[]trakt.TraktWatchedResponse
	TraktRequest *trakt.MarkAsWatchedRequest
	EmbyItems    *[]emby.EmbyItemResponse
}

func syncEmbyMovie(ctx *context.Context, cfg *config.ConfigEntity, request *SyncRequest) error {
	// Create a map with IMDb IDs as keys and TraktWatchedItem as values
	imdbMap := make(map[string]*trakt.TraktWatchedItem)

	// Iterate over TraktItems (movies) and populate the map
	for _, movie := range *request.TraktItems {
		if movie.Movie.Ids.Imdb != "" {
			imdbMap[movie.Movie.Ids.Imdb] = &movie.Movie
		}
	}

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

		_, traktExists := imdbMap[imdbId]

		fmt.Println(spacePrefix1, "Processing Emby movie:", embyItem.Name, "with IMDb ID:", imdbId)
		if !traktExists && embyItem.UserData.Played {
			fmt.Println(spacePrefix2, "Adding new movie to Trakt request")
			request.TraktRequest.Movies = append(request.TraktRequest.Movies, trakt.MarkAsWatchedMovieRequest{
				Ids: trakt.MarkAsWatchedIds{
					Imdb: imdbId,
				},
				WatchedAt: embyItem.UserData.LastPlayedDate,
			})
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

	// Create a map with IMDb IDs as keys and TraktWatchedItem as values
	traktImdbMap := make(map[string]map[int16]map[int16]*trakt.TraktWatchedResponse)

	fmt.Println("== Processing Emby items for TV shows ==")
	for _, show := range *request.TraktItems {
		if show.Show.Ids.Imdb != "" {
			imdbId := show.Show.Ids.Imdb
			if _, exists := traktImdbMap[imdbId]; !exists {
				traktImdbMap[imdbId] = make(map[int16]map[int16]*trakt.TraktWatchedResponse)
			}

			for _, season := range show.Seasons {
				seasonNumber := season.Number
				if _, exists := traktImdbMap[imdbId][seasonNumber]; !exists {
					traktImdbMap[imdbId][seasonNumber] = make(map[int16]*trakt.TraktWatchedResponse)
				}

				for _, episode := range season.Episodes {
					episodeNumber := episode.Number
					traktImdbMap[imdbId][seasonNumber][episodeNumber] = &show
				}
			}
		}
	}

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

			_, traktImdbEpisodeExists := traktImdbMap[imdbId][embySeasonNumber][embyEpisodeNumber]
			if !traktImdbEpisodeExists && embyShow.UserData.Played {
				fmt.Printf("%s Marking Trakt episode as watched for S%dE%d\n", spacePrefix2, embySeasonNumber, embyEpisodeNumber)

				foundShow, found := utils.FindBy(&request.TraktRequest.Shows, func(item trakt.MarkAsWatchedShowRequest) bool {
					return item.Ids.Imdb == imdbId
				})

				if found {
					foundSeason, found := utils.FindBy(&foundShow.Seasons, func(item trakt.MarkAsWatchedSeasonsRequest) bool {
						return item.Number == embySeasonNumber
					})
					if found {
						foundSeason.Episodes = append(foundSeason.Episodes, trakt.MarkAsWatchedEpisodes{
							Number:    embyEpisodeNumber,
							WatchedAt: embyShow.UserData.LastPlayedDate,
						})
					} else {
						foundShow.Seasons = append(foundShow.Seasons, trakt.MarkAsWatchedSeasonsRequest{
							WatchedAt: embyShow.UserData.LastPlayedDate,
							Number:    embySeasonNumber,
							Episodes: []trakt.MarkAsWatchedEpisodes{
								{
									Number:    embyEpisodeNumber,
									WatchedAt: embyShow.UserData.LastPlayedDate,
								},
							},
						})
					}
				} else {
					request.TraktRequest.Shows = append(request.TraktRequest.Shows, trakt.MarkAsWatchedShowRequest{
						Ids: trakt.MarkAsWatchedIds{
							Imdb: imdbId,
						},
						Seasons: []trakt.MarkAsWatchedSeasonsRequest{
							{
								WatchedAt: embyShow.UserData.LastPlayedDate,
								Number:    embySeasonNumber,
								Episodes: []trakt.MarkAsWatchedEpisodes{
									{
										Number:    embyEpisodeNumber,
										WatchedAt: embyShow.UserData.LastPlayedDate,
									},
								},
							},
						},
					})
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
