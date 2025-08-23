using System.Net.Http.Json;
using System.Reflection;
using Microsoft.AspNetCore.WebUtilities;
using Microsoft.Extensions.Logging;
using TraktSync.Config;
using TraktSync.Emby.Models;

namespace TraktSync.Emby;

public class EmbyClient(
    ConfigHandler configHandler,
    ILogger<EmbyClient> logger)
{
    public async Task<EmbyResponse> GetTvShowsSync()
    {
        var tvShows = await GetItemsAsync(EmbyItemType.Series);
        var episodes = await GetItemsAsync(EmbyItemType.Episode);
        var tvSeriesDic = tvShows.Items?.ToDictionary(x => x.Id, x => x);
        foreach (var episode in episodes?.Items?.Where(x => !string.IsNullOrEmpty(x.ParentId))!)
        {
            if (tvSeriesDic != null && tvSeriesDic.TryGetValue(episode.ParentId!, out var tvShow))
            {
                tvShow.Episodes ??= new List<EmbyItemResponse>();
                tvShow.Episodes.Add(episode);
            }
            else
            {
                logger.LogError(
                    "Emby client | Episode {EpisodeId} has no matching TV show with ID {ParentId}",
                    episode.Id, episode.ParentId);
            }
        }

        return tvShows;
    }

    public async Task<EmbyResponse> GetMoviesSync() =>
        await GetItemsAsync(EmbyItemType.Movie);
    
    private async Task<EmbyResponse> GetItemsAsync(EmbyItemType type)
    {
        var config = configHandler.GetAsync()?.Emby ?? throw new NullReferenceException("Emby config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("X-Emby-Token", config.ApiKey);

        var query = new EmbyRequest
        {
            IncludeItemTypes = type.ToString()
        };
        var queryParams = query
            .GetType()
            .GetProperties(BindingFlags.Instance | BindingFlags.Public)
            .ToDictionary(
                prop => prop.Name,
                prop => prop.GetValue(query, null)?.ToString()
            );

        var url = QueryHelpers.AddQueryString($"/Users/{config.UserId}/Items", queryParams);

        try
        {
            var response = await httpClient.GetAsync(url);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Emby client | Error getting watched shows: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Emby client | Error getting watched shows");    
            }
            
            var result = await response.Content.ReadFromJsonAsync<EmbyResponse>();
            return result!;
        }
        catch (Exception ex)
        {
            logger.LogError("Emby client | Error getting watched shows: {RequestMessage}", ex.Message);

            throw;
        }
    }
    
    public async Task MarkAsWatchedAsync(string itemId)
    {
        var config = configHandler.GetAsync()?.Emby ?? throw new NullReferenceException("Emby config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("X-Emby-Token", config.ApiKey);
        
        try
        {
            var url = $"/Users/{config.UserId}/PlayedItems/{itemId}";
            var response = await httpClient.PostAsJsonAsync(url, new {});
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Emby client | Error marking as watched item: {StatusCode} - {RequestMessage} - {ItemId}",
                    response.StatusCode, response.RequestMessage, itemId);
                throw new Exception("Emby client | Error making as watched shows");    
            }
        }
        catch (Exception ex)
        {
            logger.LogError("Emby client | Error getting watched shows: {RequestMessage}", ex.Message);

            throw;
        }
    }
}