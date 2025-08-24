using Microsoft.Extensions.Logging;
using TraktSync.Emby;
using TraktSync.Trakt;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler;

public class SyncMoviesHandler(
    TraktClient traktClient,
    EmbyClient embyClient,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync(TraktMarkAsWatchedRequest traktRequest)
    {
        logger.LogInformation("Sync movies | Starting sync process");
        
        var traktWatchedMovies = await traktClient.GetWatchedMoviesAsync();
        var embyMovies = await embyClient.GetMoviesSync();
        
        var traktMoviesDic = traktWatchedMovies
            .Where(w => !string.IsNullOrEmpty(w.Movie?.Ids?.Imdb))
            .GroupBy(x => x.Movie?.Ids?.Imdb ?? string.Empty)
            .ToDictionary(x => x.Key, x => x.Last());

        traktRequest.Movies ??= [];
        
        foreach (var embyMovie in embyMovies.Items ?? [])
        {
            if (traktMoviesDic.TryGetValue(embyMovie.Ids?.Imdb ?? string.Empty, out var traktMovie))
            {
                if (embyMovie.Data?.Played == true)
                {
                    // do nothing, already marked as watched
                }
                else
                {
                    // mark as watched in Emby
                    await embyClient.MarkAsWatchedAsync(embyMovie.Id);
                    logger.LogInformation("Sync movies | Marked movie {Movie} as watched on emby", embyMovie.Name);
                }
            }
            else
            {
                if (embyMovie.Data?.Played == true)
                {
                    // mark as watched in Trakt
                    traktRequest.Movies.Add(new TraktMarkAsWatchedMovieRequest
                    {
                        Ids = new TraktMarkAsWatchedIdsRequest { Imdb = embyMovie.Ids?.Imdb },
                        WatchedAt = embyMovie.Data?.LastPlayedDate ?? DateTime.UtcNow
                    });
                    logger.LogInformation("Sync movies | Marked movie {Movie} as watched on trakt", embyMovie.Name);
                }
            }
        }
        
        logger.LogInformation("Sync movies | Sync process completed");
    }
}