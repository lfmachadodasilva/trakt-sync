using Microsoft.Extensions.Logging;
using TraktSync.Trakt;

namespace TraktSync.Handler;

public class SyncMoviesHandler(TraktClient traktClient, ILogger<SyncHandler> logger)
{
    public async Task SyncAsync()
    {
        logger.LogInformation("Sync movies | Starting sync process");
        
        var traktWatchedMovies = await traktClient.GetWatchedMoviesAsync();
        
        logger.LogInformation("Sync movies | Sync process completed");
    }
}