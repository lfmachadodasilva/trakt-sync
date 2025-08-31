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
    public async Task SyncAsync(CancellationToken cancellationToken = default)
    {
        logger.LogInformation("Sync | Starting sync process");

        TraktMarkAsWatchedRequest traktRequest = new();
        
        await syncMoviesHandler.SyncAsync(traktRequest, cancellationToken);
        await syncTvShowsHandler.SyncAsync(traktRequest, cancellationToken);
        await traktClient.MarkAsWatchedAsync(traktRequest, cancellationToken: cancellationToken);

        logger.LogInformation("Sync | Sync process completed");
    }
}