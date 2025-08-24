using System.Net.Http.Json;
using System.Reflection;
using System.Text.Json.Serialization;
using Microsoft.AspNetCore.WebUtilities;
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

    public async Task<string> CodeAsync()
    {
        var config = configHandler.GetAsync()?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;

        var authCode = new TraktAuthCode
        {
            ClientId = config.ClientId,
            RedirectUrl = config.RedirectUrl
        };
        
        var queryParams = authCode
            .GetType()
            .GetProperties(BindingFlags.Instance | BindingFlags.Public)
            .ToDictionary(
                prop => prop.GetCustomAttribute<JsonPropertyNameAttribute>()?.Name ?? prop.Name,
                prop => prop.GetValue(authCode, null)?.ToString()
            );
        var url = QueryHelpers.AddQueryString("oauth/authorize", queryParams);
        return config.BaseUrl + url;
    }
    
    public async Task<TraktAuthResponse> AuthAsync(string code)
    {
        var config = configHandler.GetAsync()?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");

        var request = new TraktAuthRequest
        {
            Code = code,
            ClientId = config.ClientId,
            ClientSecret = config.ClientSecret,
            RedirectUrl = config.RedirectUrl,
            GrantType = "authorization_code"
        };
        
        try
        {
            var response = await httpClient.PostAsJsonAsync("/oauth/token", request);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Trakt client | Error marking as watched: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Trakt client | Error marking as watched");    
            }
            var result = await response.Content.ReadFromJsonAsync<TraktAuthResponse>();
            return result!;
        }
        catch (Exception ex)
        {
            logger.LogError("Trakt client | Error marking as watched: {RequestMessage}", ex.Message);
            throw;
        }
    }
    
    public async Task<TraktAuthResponse> AuthRefreshAccessTokenAsync()
    {
        var config = configHandler.GetAsync()?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("trakt-api-version", config.ApiVersion);
        httpClient.DefaultRequestHeaders.Add("trakt-api-key", config.ClientId);
        httpClient.DefaultRequestHeaders.Add("Authorization", config.AccessToken);
        
        var request = new TraktAuthRefreshRequest
        {
            ClientId = config.ClientId,
            ClientSecret = config.ClientSecret,
            RefreshToken = config.RefreshToken,
            RedirectUrl = config.RedirectUrl,
            GrantType = "refresh_token"
        };
        
        try
        {
            var response = await httpClient.PostAsJsonAsync("/oauth/token", request);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Trakt client | Error marking as watched: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Trakt client | Error marking as watched");    
            }
            
            var result = await response.Content.ReadFromJsonAsync<TraktAuthResponse>();
            return result!;
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