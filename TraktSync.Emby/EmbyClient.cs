using System.Net.Http.Json;
using System.Reflection;
using Microsoft.AspNetCore.WebUtilities;
using Microsoft.Extensions.Logging;
using TraktSync.Config;
using TraktSync.Emby.Models;

namespace TraktSync.Emby;

public interface IEmbyClient
{
    Task<EmbyResponse?> GetTvShowsSync(CancellationToken cancellationToken = default);
    Task<EmbyResponse?> GetMoviesSync(CancellationToken cancellationToken = default);
    Task MarkAsWatchedAsync(string itemId, CancellationToken cancellationToken = default);
}

public class EmbyClient(
    ConfigHandler configHandler,
    ILogger<EmbyClient> logger) : IEmbyClient
{
    public async Task<EmbyResponse?> GetTvShowsSync(CancellationToken cancellationToken = default)
    {
        var tvShows = await GetItemsAsync(EmbyItemType.Series);
        var episodes = await GetItemsAsync(EmbyItemType.Episode);
        var tvSeriesDic = tvShows?.Items?.ToDictionary(x => x.Id, x => x);
        foreach (var episode in episodes?.Items?.Where(x => !string.IsNullOrEmpty(x.ParentId))!)
        {
            if (tvSeriesDic != null && tvSeriesDic.TryGetValue(episode.ParentId!, out var tvShow))
            {
                tvShow.Episodes ??= new List<EmbyItemResponse>();
                episode.Name = $"{tvShow.Name}|{episode.Name}|S{episode.Season:00}E{episode.Episode:00}";
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

    public async Task<EmbyResponse?> GetMoviesSync(CancellationToken cancellationToken = default) =>
        await GetItemsAsync(EmbyItemType.Movie, cancellationToken);
    
    private async Task<EmbyResponse?> GetItemsAsync(EmbyItemType type, CancellationToken cancellationToken = default)
    {
        try
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
            
            var response = await httpClient.GetAsync(url, cancellationToken);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Emby client | Error getting watched shows: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Emby client | Error getting watched shows");    
            }
            
            var result = await response.Content.ReadFromJsonAsync<EmbyResponse>(cancellationToken);
            return result!;
        }
        catch (Exception ex)
        {
            logger.LogError("Emby client | Error getting watched shows: {RequestMessage}", ex.Message);
            return null;
        }
    }
    
    public async Task MarkAsWatchedAsync(string itemId, CancellationToken cancellationToken = default)
    {
        var config = configHandler.GetAsync()?.Emby ?? throw new NullReferenceException("Emby config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("X-Emby-Token", config.ApiKey);
        
        try
        {
            var url = $"/Users/{config.UserId}/PlayedItems/{itemId}";
            var response = await httpClient.PostAsJsonAsync(url, new {}, cancellationToken: cancellationToken);
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