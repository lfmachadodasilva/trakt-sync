using Microsoft.Extensions.DependencyInjection;
using TraktSync.Trakt.Models;

namespace TraktSync.Trakt;

public static class Extensions
{
    public static IServiceCollection AddTrakt(this IServiceCollection services) =>
        services.AddSingleton<ITraktClient, TraktClient>();
    
    public static TraktMarkAsWatchedRequest AddMarkAsWatchedRequest(
        this TraktMarkAsWatchedRequest request, string imdb, short season, short episode, DateTime watchedAt)
    {
        request.TvShows ??= new List<TraktMarkAsWatchedTvShowRequest>();
        
        var showRequest = request.TvShows.FirstOrDefault(s => s.Ids?.Imdb == imdb);
        if (showRequest is null)
        {
            showRequest = new TraktMarkAsWatchedTvShowRequest
            {
                Ids = new TraktMarkAsWatchedIdsRequest { Imdb = imdb },
                Seasons = new List<TraktMarkAsWatchedSeasonRequest>()
            };
            request.TvShows.Add(showRequest);
        }
        
        var seasonRequest = showRequest.Seasons?.FirstOrDefault(s => s.Number == season);
        if (seasonRequest is null)
        {
            seasonRequest = new TraktMarkAsWatchedSeasonRequest
            {
                Number = season,
                Episodes = new List<TraktMarkAsWatchedEpisodeRequest>
                {
                    new()
                    {
                        Number = episode,
                        WatchedAt = watchedAt
                    }
                }
            };
            showRequest.Seasons?.Add(seasonRequest);
        }
        
        var episodeRequest = seasonRequest.Episodes?.FirstOrDefault(e => e.Number == episode);
        if (episodeRequest is null)
        {
            episodeRequest = new TraktMarkAsWatchedEpisodeRequest
            {
                Number = episode,
                WatchedAt = watchedAt
            };
            seasonRequest.Episodes?.Add(episodeRequest);
        }

        return request;
    }
}