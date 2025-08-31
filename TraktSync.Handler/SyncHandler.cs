using Microsoft.Extensions.Logging;
using TraktSync.Trakt;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler;

public class SyncHandler(
    SyncTvShowsHandler syncTvShowsHandler, 
    SyncMoviesHandler syncMoviesHandler,
    ITraktClient traktClient,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync()
    {
        logger.LogInformation("Sync | Starting sync process");

        TraktMarkAsWatchedRequest traktRequest = new();
        
        await syncMoviesHandler.SyncAsync(traktRequest);
        await syncTvShowsHandler.SyncAsync(traktRequest);
        await traktClient.MarkAsWatchedAsync(traktRequest);

        logger.LogInformation("Sync | Sync process completed");
    }
}