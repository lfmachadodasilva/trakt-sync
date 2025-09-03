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
    Task<ICollection<TraktWatchedTvShowResponse>> GetWatchedTvShowsAsync(CancellationToken cancellationToken = default);
    Task<ICollection<TraktWatchedMoviesResponse>> GetWatchedMoviesAsync(CancellationToken cancellationToken = default);
    Task MarkAsWatchedAsync(
        TraktMarkAsWatchedRequest traktRequest,
        bool refreshToken = true,
        CancellationToken cancellationToken = default);
    string GetCodeUrl();
    Task AuthAsync(string code, CancellationToken cancellationToken = default);
    Task AuthRefreshAccessTokenAsync(CancellationToken cancellationToken = default);
}
    
public class TraktClient(
    ConfigHandler configHandler,
    ILogger<TraktClient> logger) : ITraktClient
{
    public async Task<ICollection<TraktWatchedTvShowResponse>> GetWatchedTvShowsAsync(CancellationToken cancellationToken = default) =>
        await GetWatchedAsync<TraktWatchedTvShowResponse>(TraktWatchedType.Shows, cancellationToken);
    
    public async Task<ICollection<TraktWatchedMoviesResponse>> GetWatchedMoviesAsync(CancellationToken cancellationToken = default) =>
        await GetWatchedAsync<TraktWatchedMoviesResponse>(TraktWatchedType.Movies, cancellationToken);
    
    public async Task MarkAsWatchedAsync(
        TraktMarkAsWatchedRequest traktRequest,
        bool refreshToken = true,
        CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(traktRequest);
            
        var config = configHandler.GetAsync()?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("trakt-api-version", config.ApiVersion);
        httpClient.DefaultRequestHeaders.Add("trakt-api-key", config.ClientId);
        httpClient.DefaultRequestHeaders.Add("Authorization", $"{config.TokenType} {config.AccessToken}");
        
        logger.LogInformation("Trakt client | Marking as watched: {Config}",
            System.Text.Json.JsonSerializer.Serialize(config));
        
        try
        {
            var response = await httpClient.PostAsJsonAsync("sync/history", traktRequest, cancellationToken);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Trakt client | Error marking as watched: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                
                if (refreshToken)
                {
                    await AuthRefreshAccessTokenAsync(cancellationToken);
                    await MarkAsWatchedAsync(traktRequest, false, cancellationToken);
                    return;
                }

                throw new Exception("Trakt client | Error marking as watched");    
            }
        }
        catch (Exception ex)
        {
            logger.LogError("Trakt client | Error marking as watched: {RequestMessage}", ex.Message);
            throw;
        }
    }

    public string GetCodeUrl()
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
    
    public async Task AuthAsync(string code, CancellationToken cancellationToken = default)
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
            var response = await httpClient.PostAsJsonAsync("/oauth/token", request, cancellationToken);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Trakt client | Error marking as watched: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Trakt client | Error marking as watched");    
            }
            var result = await response.Content.ReadFromJsonAsync<TraktAuthResponse>(cancellationToken);
            
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
    
    public async Task AuthRefreshAccessTokenAsync(CancellationToken cancellationToken = default)
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
            var response = await httpClient.PostAsJsonAsync("/oauth/token", request, cancellationToken);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Trakt client | Error refreshing trakt token: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Trakt client | Error refreshing trakt token");
            }
            
            var result = await response.Content.ReadFromJsonAsync<TraktAuthResponse>(cancellationToken);
            
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
    
    private async Task<ICollection<T>> GetWatchedAsync<T>(
        TraktWatchedType type,
        CancellationToken cancellationToken = default)
    {
        var config = configHandler.GetAsync()?.Trakt ?? throw new NullReferenceException("Trakt config is null");
        
        using HttpClient httpClient = new();
        httpClient.BaseAddress = config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("trakt-api-version", config.ApiVersion);
        httpClient.DefaultRequestHeaders.Add("trakt-api-key", config.ClientId);

        try
        {
            var response = await httpClient.GetAsync($"/users/lfmachadodasilva/watched/{type.ToString().ToLower()}", cancellationToken);
            if (!response.IsSuccessStatusCode)
            {
                logger.LogError(
                    "Trakt client | Error getting watched shows: {StatusCode} - {RequestMessage}",
                    response.StatusCode, response.RequestMessage);
                throw new Exception("Trakt client | Error getting watched shows");    
            }
            
            var result = await response.Content.ReadFromJsonAsync<ICollection<T>>(cancellationToken);
            return result!;
        }
        catch (Exception ex)
        {
            logger.LogError("Trakt client | Error getting watched shows: {RequestMessage}", ex.Message);
            throw;
        }
    }
}