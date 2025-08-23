using Microsoft.Extensions.Logging;
using TraktSync.Emby;
using TraktSync.Trakt;

namespace TraktSync.Handler;

public class SyncMoviesHandler(
    TraktClient traktClient,
    EmbyClient embyClient,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync()
    {
        logger.LogInformation("Sync movies | Starting sync process");
        
        var traktWatchedMovies = await traktClient.GetWatchedMoviesAsync();
        var embyMovies = await embyClient.GetMovies();
        
        logger.LogInformation("Sync movies | Sync process completed");
    }
}