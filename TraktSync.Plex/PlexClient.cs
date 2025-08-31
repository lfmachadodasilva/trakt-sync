using System.Net.Http.Json;
using LukeHagar.PlexAPI.SDK;
using LukeHagar.PlexAPI.SDK.Models.Requests;
using Microsoft.Extensions.Logging;
using TraktSync.Config;
using TraktSync.Plex.Models;

namespace TraktSync.Plex;

public interface IPlexClient
{
    Task<ICollection<PlexTvShow>?> GetTvShowsSync(CancellationToken cancellationToken = default);
    Task<GetLibrarySectionsAllResponse?> GetMoviesSync(CancellationToken cancellationToken = default);
    Task MarkAsWatchedAsync(string itemId, CancellationToken cancellationToken = default);
}

public class PlexClient(
    ConfigHandler configHandler,
    ILogger<PlexClient> logger) : IPlexClient
{
    public async Task<GetLibrarySectionsAllResponse?> GetMoviesSync(CancellationToken cancellationToken = default)
    {
        try
        {
            var config = configHandler.GetAsync()?.Plex ?? throw new NullReferenceException("Plex config is null");
        
            using var client = new SpeakeasyHttpClient2();
            var sdk = new PlexAPI(accessToken: config.ApiKey, serverUrl: config.BaseUrl.ToString(), client: client);
            var request = new GetLibrarySectionsAllRequest
            {
                SectionKey = 1,
                Type = GetLibrarySectionsAllQueryParamType.Movie,
                IncludeGuids = QueryParamIncludeGuids.Enable,
                IncludeAdvanced = IncludeAdvanced.Enable,
                IncludeCollections = QueryParamIncludeCollections.Enable,
                IncludeExternalMedia = QueryParamIncludeExternalMedia.Disable,
                IncludeMeta = GetLibrarySectionsAllQueryParamIncludeMeta.Enable
            };
            
            var result = await sdk.Library.GetLibrarySectionsAllAsync(request);
            return result;
        }
        catch (Exception ex)
        {
            logger.LogError(ex, "Plex client | Error getting watched movies");
            return null;
        }
    }
    
    public async Task<ICollection<PlexTvShow>?> GetTvShowsSync(CancellationToken cancellationToken = default)
    {
        try
        {
            var config = configHandler.GetAsync()?.Plex ?? throw new NullReferenceException("Plex config is null");
        
            using var client = new SpeakeasyHttpClient2();
            var sdk = new PlexAPI(accessToken: config.ApiKey, serverUrl: config.BaseUrl.ToString(), client: client);
            var request = new GetLibrarySectionsAllRequest
            {
                SectionKey = 2,
                Type = GetLibrarySectionsAllQueryParamType.Episode,
                IncludeGuids = QueryParamIncludeGuids.Enable,
                IncludeAdvanced = IncludeAdvanced.Enable,
                IncludeCollections = QueryParamIncludeCollections.Enable,
                IncludeExternalMedia = QueryParamIncludeExternalMedia.Disable,
                IncludeMeta = GetLibrarySectionsAllQueryParamIncludeMeta.Enable
            };
            
            ICollection<PlexTvShow> tvShows = [];
            var response = await sdk.Library.GetLibrarySectionsAllAsync(request);
            
            foreach (var tvShow in response.Object?.MediaContainer?.Metadata ?? [])
            {
                var plexTvShow = new PlexTvShow
                {
                    Id = tvShow.RatingKey,
                    Type = tvShow.Type,
                    Imdb = tvShow.Guids?.Select(x => x.Id)?.GetImdb(),
                    Name = tvShow.Title,
                    Object = tvShow
                };
                tvShows.Add(plexTvShow);
                
                var seasons = await GetChildrenAsync(tvShow.RatingKey ?? string.Empty, cancellationToken);

                foreach (var season in seasons?.MediaContainer?.Metadata ?? [])
                {
                    var plexSeason = new PlexTvShow
                    {
                        Id = season?.RatingKey,
                        Type = season?.Type,
                        Season = season?.Index ?? 0,
                        Name = season?.Title,
                        Object = season
                    };
                    plexTvShow.Children.Add(plexSeason);
                    
                    var episodes = await GetChildrenAsync(season?.RatingKey ?? string.Empty, cancellationToken);
                    foreach (var episode in episodes?.MediaContainer?.Metadata ?? [])
                    {
                        var plexEpisode = new PlexTvShow
                        {
                            Id = episode?.RatingKey,
                            Type = episode?.Type,
                            Season = episode?.ParentIndex ?? 0,
                            Episode = episode?.Index ?? 0,
                            Object = episode,
                            PlayedAt = episode?.LastViewedAt is not null ?
                                DateTimeOffset.FromUnixTimeSeconds(episode.LastViewedAt ?? 0).UtcDateTime : null,
                            Played = (episode?.ViewCount ?? 0) > 0
                        };
                        plexEpisode.Name = $"{plexTvShow.Name}|{episode?.Title ?? string.Empty}|S{plexEpisode.Season:00}E{plexEpisode.Episode:00}";
                        plexSeason.Children.Add(plexEpisode);
                    }
                }
            }
            
            return tvShows;
        }
        catch (Exception ex)
        {
            logger.LogError(ex, "Plex client | Error getting watched movies");
            return null;
        }
    }
    
    private async Task<GetMetadataChildrenResponseBody> GetChildrenAsync(
        string plexId,
        CancellationToken cancellationToken = default)
    {
        var config = configHandler.GetAsync()?.Plex ?? throw new NullReferenceException("Plex config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("X-Plex-Token", config.ApiKey);

        try
        {
            var response = await httpClient.GetAsync($"/library/metadata/{plexId}/children", cancellationToken);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Emby client | Error getting watched shows: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Emby client | Error getting watched shows");    
            }
            
            var result = await response.Content.ReadFromJsonAsync<GetMetadataChildrenResponseBody>(cancellationToken);
            return result!;
        }
        catch (Exception ex)
        {
            logger.LogError("Plex client | Error getting watched movies: {RequestMessage}", ex.Message);
            throw;
        }
    }
    
    public async Task MarkAsWatchedAsync(string itemId, CancellationToken cancellationToken = default)
    {
        var config = configHandler.GetAsync()?.Plex ?? throw new NullReferenceException("Plex config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("X-Plex-Token", config.ApiKey);
        
        // using var client = new SpeakeasyHttpClient2();
        // var sdk = new PlexAPI(accessToken: config.ApiKey, serverUrl: config.BaseUrl.ToString(), client: client);
        
        try
        {
            // await sdk.Media.MarkPlayedAsync(key: itemId);
            var response = await httpClient.GetAsync(
                $"/:/scrobble?key={itemId}&identifier=com.plexapp.plugins.library",
                cancellationToken);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Emby client | Error getting watched shows: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Emby client | Error getting watched shows");    
            }
        }
        catch (Exception ex)
        {
            logger.LogError("Plex client | Error getting watched shows: {RequestMessage}", ex.Message);
            throw;
        }
    }
}