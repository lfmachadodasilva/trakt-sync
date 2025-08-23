using Microsoft.Extensions.Logging;
using TraktSync.Emby;
using TraktSync.Trakt;

namespace TraktSync.Handler;

public class SyncTvShowsHandler(
    TraktClient traktClient,
    EmbyClient embyClient,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync()
    {
        logger.LogInformation("Sync tv shows | Starting sync process");
        
        var traktWatchedTvShows = await traktClient.GetWatchedTvShowsAsync();
        var embyTvShows = await embyClient.GetTvShowsSync();
        
        logger.LogInformation("Sync tv shows | Sync process completed");
    }
}