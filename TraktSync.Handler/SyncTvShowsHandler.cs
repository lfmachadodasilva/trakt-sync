using Microsoft.Extensions.Logging;
using TraktSync.Trakt;

namespace TraktSync.Handler;

public class SyncTvShowsHandler(TraktClient traktClient, ILogger<SyncHandler> logger)
{
    public async Task SyncAsync()
    {
        logger.LogInformation("Sync tv shows | Starting sync process");
        
        var traktWatchedTvShows = await traktClient.GetWatchedTvShowsAsync();
        
        logger.LogInformation("Sync tv shows | Sync process completed");
    }
}