using Microsoft.Extensions.Logging;
using TraktSync.Emby;
using TraktSync.Plex;
using TraktSync.Trakt;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler;

using TraktTvShowsDictionary = Dictionary<string, Dictionary<short, HashSet<short>>>;

public class SyncTvShowsHandler(
    ITraktClient traktClient,
    IEmbyClient embyClient,
    IPlexClient plexClient,
    ILogger<SyncHandler> logger)
{
    public async Task SyncAsync(TraktMarkAsWatchedRequest traktRequest)
    {
        ArgumentNullException.ThrowIfNull(traktRequest, nameof(traktRequest));
        
        logger.LogInformation("Sync tv shows | Starting sync process");
        
        var traktWatchedTvShows = await traktClient.GetWatchedTvShowsAsync();
        var traktTvShowsDic = ToDictionary(traktWatchedTvShows);
        
        await SyncEmbyAsync(traktRequest, traktTvShowsDic);
        await SyncPlexAsync(traktRequest, traktTvShowsDic);
        
        logger.LogInformation("Sync tv shows | Sync process completed");
    }

    private async Task SyncEmbyAsync(
        TraktMarkAsWatchedRequest traktRequest,
        TraktTvShowsDictionary traktTvShowsDic)
    {
        var embyTvShows = await embyClient.GetTvShowsSync();
        
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
    }

    private async Task SyncPlexAsync(
        TraktMarkAsWatchedRequest traktRequest,
        TraktTvShowsDictionary traktTvShowsDic)
    {
        var plexTvShows = await plexClient.GetTvShowsSync();
        
        foreach (var tvShow in plexTvShows ?? [])
        {
            var imdb = tvShow.Imdb ?? string.Empty;
            foreach (var seasons in tvShow.Children ?? [])
            {
                foreach (var episode in seasons.Children ?? [])
                {
                    var playedPlex = episode.Played ?? false;
                    var playedTrakt = traktTvShowsDic.TryGetValue(imdb ?? string.Empty, out _);
                    
                    if (playedTrakt && !playedPlex)
                    {
                        // mark as watched in plex
                        await plexClient.MarkAsWatchedAsync(episode.Id ?? string.Empty);
                    }
                    else if (!playedTrakt && playedPlex)
                    {
                        traktRequest.AddMarkAsWatchedRequest(
                            imdb ?? string.Empty,
                            (short) (episode.Season ?? 0),
                            (short) (episode.Episode ?? 0),
                            episode.PlayedAt ?? DateTime.Now);
                    }
                }
            }
        }
    }

    private static TraktTvShowsDictionary ToDictionary(ICollection<TraktWatchedTvShowResponse> traktWatchedTvShows)
    {
        var traktTvShowsDic = new TraktTvShowsDictionary();
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
        return traktTvShowsDic;
    }
}