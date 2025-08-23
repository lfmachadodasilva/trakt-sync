using System.Net.Http.Json;
using System.Reflection;
using Microsoft.AspNetCore.WebUtilities;
using Microsoft.Extensions.Logging;
using TraktSync.Emby.Models;

namespace TraktSync.Emby;

public class EmbyClient(ILogger<EmbyClient> logger)
{
    private readonly EmbyConfig _config = new()
    {
        BaseUrl = new Uri("http://192.168.1.13:8096"),
        ApiKey = "b039ba2b065e4ba1bca2307cce593478",
        UserId = "aac3a78d9f184ea480fb1629e76aad57"
    };

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
        }

        return tvShows;
    }
        
    
    public async Task<EmbyResponse> GetMoviesSync() =>
        await GetItemsAsync(EmbyItemType.Movie);
    
    private async Task<EmbyResponse> GetItemsAsync(EmbyItemType type)
    {
        using HttpClient httpClient = new();
        httpClient.BaseAddress = _config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("X-Emby-Token", _config.ApiKey);

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

        var url = QueryHelpers.AddQueryString($"/Users/{_config.UserId}/Items", queryParams);

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
}