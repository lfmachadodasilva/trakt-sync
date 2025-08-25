using Microsoft.Extensions.Logging;
using TraktSync.Emby;
using TraktSync.Plex;
using TraktSync.Trakt;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler;

public class SyncMoviesHandler(
    ITraktClient traktClient,
    IEmbyClient embyClient,
    IPlexClient plexClient,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync(TraktMarkAsWatchedRequest traktRequest)
    {
        ArgumentNullException.ThrowIfNull(traktRequest, nameof(traktRequest));
        
        logger.LogInformation("Sync movies | Starting sync process");
        
        var traktWatchedMovies = await traktClient.GetWatchedMoviesAsync();
        var traktMoviesDic = traktWatchedMovies
            .Where(w => !string.IsNullOrEmpty(w.Movie?.Ids?.Imdb))
            .GroupBy(x => x.Movie?.Ids?.Imdb ?? string.Empty)
            .ToDictionary(x => x.Key, x => x.Last());
        traktRequest.Movies ??= [];

        await SyncPlexAsync(traktRequest, traktMoviesDic);
        await SyncEmbyAsync(traktRequest, traktMoviesDic);
        
        logger.LogInformation("Sync movies | Sync process completed");
    }
    
    private async Task SyncEmbyAsync(
        TraktMarkAsWatchedRequest traktRequest,
        Dictionary<string,TraktWatchedMoviesResponse> traktMoviesDic)
    {
        var embyMovies = await embyClient.GetMoviesSync();

        foreach (var embyMovie in embyMovies.Items ?? [])
        {
            var imdb = embyMovie.Ids?.Imdb ?? string.Empty;
            var playedEmby = embyMovie.Data?.Played == true;
            var playedTrakt = traktMoviesDic.TryGetValue(imdb, out var traktMovie);

            if (playedTrakt && !playedEmby)
            {
                // mark as watched in Emby
                await embyClient.MarkAsWatchedAsync(embyMovie.Id);
                logger.LogInformation("Sync movies | Marked movie {Movie} as watched on emby",
                    traktMovie?.Movie?.Title);
            }
            else if (!playedTrakt && playedEmby)
            {
                // mark as watched in Trakt
                traktRequest.Movies.Add(new TraktMarkAsWatchedMovieRequest
                {
                    Ids = new TraktMarkAsWatchedIdsRequest {Imdb = embyMovie.Ids?.Imdb},
                    WatchedAt = embyMovie.Data?.LastPlayedDate ?? DateTime.UtcNow
                });
                logger.LogInformation("Sync movies | Marked movie {Movie} as watched on trakt", embyMovie.Name);
            }
        }
    }

    private async Task SyncPlexAsync(
        TraktMarkAsWatchedRequest traktRequest,
        Dictionary<string,TraktWatchedMoviesResponse> traktMoviesDic)
    {
        var plexMovies = await plexClient.GetMoviesSync();

        foreach (var plexMovie in plexMovies?.Object?.MediaContainer?.Metadata ?? [])
        {
            var imdb = plexMovie?.Guids?.Select(x => x.Id)?.GetImdb();
            var playedPlex = plexMovie?.ViewCount > 0;
            var playedTrakt = traktMoviesDic.TryGetValue(imdb ?? string.Empty, out var traktMovie);

            if (playedTrakt && !playedPlex)
            {
                // mark as watched in plex
                await plexClient.MarkAsWatchedAsync(plexMovie?.RatingKey ?? string.Empty);
                logger.LogInformation("Sync movies | Marked movie {Movie} as watched on emby", traktMovie?.Movie?.Title);
            }
            else if (!playedTrakt && playedPlex)
            {
                traktRequest.Movies.Add(new TraktMarkAsWatchedMovieRequest
                {
                    Ids = new TraktMarkAsWatchedIdsRequest { Imdb = imdb },
                    WatchedAt = DateTimeOffset.FromUnixTimeSeconds(plexMovie?.LastViewedAt ?? 0).UtcDateTime 
                });
                logger.LogInformation("Sync movies | Marked movie {Movie} as watched on plex", plexMovie?.OriginalTitle);
            }
        }
    }
}