using Microsoft.Extensions.Logging;
using TraktSync.Emby;
using TraktSync.Trakt;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler;

public class SyncTvShowsHandler(
    TraktClient traktClient,
    EmbyClient embyClient,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync(TraktMarkAsWatchedRequest traktRequest)
    {
        logger.LogInformation("Sync tv shows | Starting sync process");
        
        var traktWatchedTvShows = await traktClient.GetWatchedTvShowsAsync();
        var embyTvShows = await embyClient.GetTvShowsSync();
        
        logger.LogInformation("Sync tv shows | Sync process completed");
    }
}