using LukeHagar.PlexAPI.SDK;
using LukeHagar.PlexAPI.SDK.Models.Requests;
using Microsoft.Extensions.Logging;
using TraktSync.Config;
using TraktSync.Plex.Models;

namespace TraktSync.Plex;

public interface IPlexClient
{
    Task GetTvShowsSync();
    Task<GetLibrarySectionsAllResponse> GetMoviesSync();
    Task MarkAsWatchedAsync(string itemId);
}

public class PlexClient(
    ConfigHandler configHandler,
    ILogger<PlexClient> logger) : IPlexClient
{
    public async Task<GetLibrarySectionsAllResponse> GetMoviesSync()
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

        try
        {
            var result = await sdk.Library.GetLibrarySectionsAllAsync(request);
            return result;
        }
        catch (Exception ex)
        {
            logger.LogError("Plex client | Error getting watched movies: {RequestMessage}", ex.Message);
            throw;
        }
    }
    
    public Task GetTvShowsSync()
    {
        throw new NotImplementedException();
    }
    
    public async Task MarkAsWatchedAsync(string itemId)
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
            var response = await httpClient.GetAsync($"/:/scrobble?key={itemId}&identifier=com.plexapp.plugins.library");
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