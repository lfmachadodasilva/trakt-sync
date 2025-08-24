using Microsoft.Extensions.Logging;
using TraktSync.Emby;
using TraktSync.Trakt;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler;

public class SyncTvShowsHandler(
    ITraktClient traktClient,
    IEmbyClient embyClient,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync(TraktMarkAsWatchedRequest traktRequest)
    {
        ArgumentNullException.ThrowIfNull(traktRequest, nameof(traktRequest));
        
        logger.LogInformation("Sync tv shows | Starting sync process");
        
        var traktWatchedTvShows = await traktClient.GetWatchedTvShowsAsync();
        var embyTvShows = await embyClient.GetTvShowsSync();
        
        var traktTvShowsDic = new Dictionary<string, Dictionary<short, HashSet<short>>>();
        foreach (var tvShow in traktWatchedTvShows ?? [])
        {
            var imdbId = tvShow?.Show?.Ids?.Imdb ?? string.Empty;
            
            if (!traktTvShowsDic.TryGetValue(imdbId, out var seasonsDic))
            {
                seasonsDic = new Dictionary<short, HashSet<short>>();
                traktTvShowsDic[imdbId] = seasonsDic;
            }
            
            foreach (var season in tvShow?.Seasons ?? [])
            {
                var seasonNumber = season.Number ?? 0;
                
                if (!seasonsDic.TryGetValue(seasonNumber, out var episodesSet))
                {
                    episodesSet = [];
                    seasonsDic[seasonNumber] = episodesSet;
                }
                
                foreach (var episode in season.Episodes ?? [])
                {
                    var episodeNumber = episode.Number ?? 0;
                    episodesSet.Add(episodeNumber);
                }
            }
        }
        
        foreach (var embyTvShow in embyTvShows?.Items ?? [])
        {
            var imdb = embyTvShow.Ids?.Imdb ?? string.Empty;
            
            traktTvShowsDic.TryGetValue(imdb, out var traktSeasonDic);
            
            foreach (var episode in embyTvShow.Episodes ?? [])
            {
                if (traktSeasonDic?.GetValueOrDefault(episode.Season ?? 0)?
                        .Contains(episode.Episode ?? 0) == true)
                {
                    if (episode.Data?.Played == true)
                    {
                        // do nothing, already marked as watched
                    }
                    else
                    {
                        // mark as watched in Emby
                        await embyClient.MarkAsWatchedAsync(episode.Id);
                        logger.LogInformation(
                            "Sync tv shows | Marked episode {Episode} as watched on emby", episode.Name);
                    }
                }
                else
                {
                    if (episode.Data?.Played == true)
                    {
                        // mark as watched in Trakt
                        traktRequest.AddMarkAsWatchedRequest(
                            embyTvShow.Ids?.Imdb ?? string.Empty,
                            episode.Season ?? 0,
                            episode.Episode ?? 0,
                            episode.Data?.LastPlayedDate ?? DateTime.UtcNow);
                        logger.LogInformation(
                            "Sync tv shows | Marked episode {Episode} as watched on trakt", episode.Name);
                    }
                }
            }
        }
        
        logger.LogInformation("Sync tv shows | Sync process completed");
    }
}