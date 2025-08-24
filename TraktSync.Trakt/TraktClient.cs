using System.Net.Http.Json;
using System.Reflection;
using System.Text.Json.Serialization;
using Microsoft.AspNetCore.WebUtilities;
using Microsoft.Extensions.Logging;
using TraktSync.Config;
using TraktSync.Trakt.Models;

namespace TraktSync.Trakt;

public interface ITraktClient
{
    Task<ICollection<TraktWatchedTvShowResponse>> GetWatchedTvShowsAsync();
    Task<ICollection<TraktWatchedMoviesResponse>> GetWatchedMoviesAsync();
    Task MarkAsWatchedAsync(TraktMarkAsWatchedRequest traktRequest, bool refreshToken = true);
    Task<string> CodeAsync();
    Task AuthAsync(string code);
    Task AuthRefreshAccessTokenAsync();
}
    
public class TraktClient(
    ConfigHandler configHandler,
    ILogger<TraktClient> logger) : ITraktClient
{
    public async Task<ICollection<TraktWatchedTvShowResponse>> GetWatchedTvShowsAsync() =>
        await GetWatchedAsync<TraktWatchedTvShowResponse>(TraktWatchedType.Shows);
    
    public async Task<ICollection<TraktWatchedMoviesResponse>> GetWatchedMoviesAsync() =>
        await GetWatchedAsync<TraktWatchedMoviesResponse>(TraktWatchedType.Movies);
    
    public async Task MarkAsWatchedAsync(TraktMarkAsWatchedRequest traktRequest, bool refreshToken = true)
    {
        ArgumentNullException.ThrowIfNull(traktRequest);
            
        var config = configHandler.GetAsync()?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("trakt-api-version", config.ApiVersion);
        httpClient.DefaultRequestHeaders.Add("trakt-api-key", config.ClientId);
        httpClient.DefaultRequestHeaders.Add("Authorization", $"{config.TokenType} {config.AccessToken}");
        
        try
        {
            var response = await httpClient.PostAsJsonAsync("sync/history", traktRequest);
            if (!response.IsSuccessStatusCode)
            {
                if (refreshToken)
                {
                    await AuthRefreshAccessTokenAsync();
                    await MarkAsWatchedAsync(traktRequest, false);
                    return;
                }
                
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
    
    public async Task AuthAsync(string code)
    {
        ArgumentException.ThrowIfNullOrEmpty(code);
        
        var config = configHandler.GetAsync();
        var configTrakt = config?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = configTrakt.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("trakt-api-version", configTrakt.ApiVersion);
        httpClient.DefaultRequestHeaders.Add("trakt-api-key", configTrakt.ClientId);

        var request = new TraktAuthRequest
        {
            Code = code,
            ClientId = configTrakt.ClientId,
            ClientSecret = configTrakt.ClientSecret,
            RedirectUrl = configTrakt.RedirectUrl,
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
            
            if (result is null)
            {
                throw new NullReferenceException("Trakt auth response is null");
            }
            
            configTrakt.AccessToken = result.AccessToken;
            configTrakt.RefreshToken = result.RefreshToken;
            configTrakt.TokenType = result.TokenType;
            configHandler.UpdateConfig(config);
        }
        catch (Exception ex)
        {
            logger.LogError("Trakt client | Error marking as watched: {RequestMessage}", ex.Message);
            throw;
        }
    }
    
    public async Task AuthRefreshAccessTokenAsync()
    {
        var config = configHandler.GetAsync();
        var configTrakt = config?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = configTrakt.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("trakt-api-version", configTrakt.ApiVersion);
        httpClient.DefaultRequestHeaders.Add("trakt-api-key", configTrakt.ClientId);
        
        var request = new TraktAuthRefreshRequest
        {
            ClientId = configTrakt.ClientId,
            ClientSecret = configTrakt.ClientSecret,
            RefreshToken = configTrakt.RefreshToken,
            RedirectUrl = configTrakt.RedirectUrl,
            GrantType = "refresh_token"
        };
        
        try
        {
            var response = await httpClient.PostAsJsonAsync("/oauth/token", request);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Trakt client | Error refreshing trakt token: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Trakt client | Error refreshing trakt token");
            }
            
            var result = await response.Content.ReadFromJsonAsync<TraktAuthResponse>();
            
            if (result is null)
            {
                throw new NullReferenceException("Trakt auth response is null");
            }
            
            configTrakt.AccessToken = result.AccessToken;
            configTrakt.RefreshToken = result.RefreshToken;
            configTrakt.TokenType = result.TokenType;
            configHandler.UpdateConfig(config);
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