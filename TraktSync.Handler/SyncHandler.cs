using Microsoft.Extensions.Logging;
using TraktSync.Trakt;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler;

public class SyncHandler(
    SyncTvShowsHandler syncTvShowsHandler, 
    SyncMoviesHandler syncMoviesHandler,
    TraktClient traktClient,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync()
    {
        logger.LogInformation("Sync | Starting sync process");

        TraktMarkAsWatchedRequest traktRequest = new();
        
        // await syncTvShowsHandler.SyncAsync(traktRequest);
        await syncMoviesHandler.SyncAsync(traktRequest);
        
        logger.LogInformation("Sync | Starting trakt process");
        await traktClient.MarkAsWatchedAsync(traktRequest);
        logger.LogInformation("Sync | Sync trakt completed");

        logger.LogInformation("Sync | Sync process completed");
    }
}