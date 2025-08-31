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
        var movies = await embyClient.GetMoviesSync();

        foreach (var movie in movies.Items ?? [])
        {
            var imdb = movie.Ids?.Imdb ?? string.Empty;
            var playedEmby = movie.Data?.Played == true;
            var playedTrakt = traktMoviesDic.TryGetValue(imdb, out var traktMovie);

            if (playedTrakt && !playedEmby)
            {
                // mark as watched in Emby
                await embyClient.MarkAsWatchedAsync(movie.Id);
                logger.LogInformation("Sync movies | Marked movie {Movie} as watched on emby",
                    traktMovie?.Movie?.Title);
            }
            else if (!playedTrakt && playedEmby)
            {
                // mark as watched in Trakt
                traktRequest.Movies.Add(new TraktMarkAsWatchedMovieRequest
                {
                    Ids = new TraktMarkAsWatchedIdsRequest {Imdb = movie.Ids?.Imdb},
                    WatchedAt = movie.Data?.LastPlayedDate ?? DateTime.UtcNow
                });
                logger.LogInformation("Sync movies | Marked movie {Movie} as watched on trakt", movie.Name);
            }
        }
    }

    private async Task SyncPlexAsync(
        TraktMarkAsWatchedRequest traktRequest,
        Dictionary<string,TraktWatchedMoviesResponse> traktMoviesDic)
    {
        var movies = await plexClient.GetMoviesSync();

        foreach (var movie in movies?.Object?.MediaContainer?.Metadata ?? [])
        {
            var imdb = movie?.Guids?.Select(x => x.Id)?.GetImdb();
            var playedPlex = movie?.ViewCount > 0;
            var playedTrakt = traktMoviesDic.TryGetValue(imdb ?? string.Empty, out var traktMovie);

            if (playedTrakt && !playedPlex)
            {
                // mark as watched in plex
                await plexClient.MarkAsWatchedAsync(movie?.RatingKey ?? string.Empty);
                logger.LogInformation("Sync movies | Marked movie {Movie} as watched on emby", traktMovie?.Movie?.Title);
            }
            else if (!playedTrakt && playedPlex)
            {
                traktRequest.Movies.Add(new TraktMarkAsWatchedMovieRequest
                {
                    Ids = new TraktMarkAsWatchedIdsRequest { Imdb = imdb },
                    WatchedAt = DateTimeOffset.FromUnixTimeSeconds(movie?.LastViewedAt ?? 0).UtcDateTime 
                });
                logger.LogInformation("Sync movies | Marked movie {Movie} as watched on plex", movie?.OriginalTitle);
            }
        }
    }
}