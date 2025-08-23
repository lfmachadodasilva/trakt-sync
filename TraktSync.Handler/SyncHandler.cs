using Microsoft.Extensions.Logging;

namespace TraktSync.Handler;

public class SyncHandler(
    SyncTvShowsHandler syncTvShowsHandler, 
    SyncMoviesHandler syncMoviesHandler,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync()
    {
        logger.LogInformation("Sync | Starting sync process");

        await syncTvShowsHandler.SyncAsync();
        await syncMoviesHandler.SyncAsync();

        logger.LogInformation("Sync | Sync process completed");
    }
}