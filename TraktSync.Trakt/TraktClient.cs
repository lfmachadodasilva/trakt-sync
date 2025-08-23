using System.Net.Http.Json;
using System.Text.Json;
using Microsoft.Extensions.Logging;
using TraktSync.Config;
using TraktSync.Trakt.Models;

namespace TraktSync.Trakt;

public class TraktClient(
    ConfigHandler configHandler,
    ILogger<TraktClient> logger)
{
    public async Task<ICollection<TraktWatchedTvShowResponse>> GetWatchedTvShowsAsync() =>
        await GetWatchedAsync<TraktWatchedTvShowResponse>(TraktWatchedType.Shows);
    
    public async Task<ICollection<TraktWatchedMoviesResponse>> GetWatchedMoviesAsync() =>
        await GetWatchedAsync<TraktWatchedMoviesResponse>(TraktWatchedType.Movies);
    
    public async Task MarkAsWatchedAsync(TraktMarkAsWatchedRequest traktRequest)
    {
        var config = configHandler.GetAsync()?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("trakt-api-version", config.ApiVersion);
        httpClient.DefaultRequestHeaders.Add("trakt-api-key", config.ClientId);
        httpClient.DefaultRequestHeaders.Add("Authorization", config.AccessToken);
        
        try
        {
            var response = await httpClient.PostAsJsonAsync($"sync/history", traktRequest);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Trakt client | Error marking as watched: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Trakt client | Error marking as watched");    
            }
        }
        catch (Exception ex)
        {
            logger.LogError("Trakt client | Error marking as watched: {RequestMessage}", ex.Message);

            throw;
        }
    }
    
    private async Task<ICollection<T>> GetWatchedAsync<T>(TraktWatchedType type)
    {
        var config = configHandler.GetAsync()?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("trakt-api-version", config.ApiVersion);
        httpClient.DefaultRequestHeaders.Add("trakt-api-key", config.ClientId);

        try
        {
            var response = await httpClient.GetAsync($"/users/lfmachadodasilva/watched/{type.ToString().ToLower()}");
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Trakt client | Error getting watched shows: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Trakt client | Error getting watched shows");    
            }
            
            var result = await response.Content.ReadFromJsonAsync<ICollection<T>>();
            return result!;
        }
        catch (Exception ex)
        {
            logger.LogError("Trakt client | Error getting watched shows: {RequestMessage}", ex.Message);

            throw;
        }
    }
}