using System.Net.Http.Json;
using System.Text.Json;
using Microsoft.Extensions.Logging;
using TraktSync.Trakt.Models;

namespace TraktSync.Trakt;

public class TraktClient(ILogger<TraktClient> logger)
{
    private readonly TraktConfig _config = new()
    {
        UserName = "lfmachadodasilva",
        ClientId = "eb4ede9a384157e9aa60aad8c72c36c0485215659c82ad7b1fe965359a55caf4"
    };
    
    public async Task<ICollection<TraktWatchedTvShowResponse>> GetWatchedTvShowsAsync() =>
        await GetWatchedAsync<TraktWatchedTvShowResponse>(TraktWatchedType.Shows);
    
    public async Task<ICollection<TraktWatchedMoviesResponse>> GetWatchedMoviesAsync() =>
        await GetWatchedAsync<TraktWatchedMoviesResponse>(TraktWatchedType.Movies);
    
    public Task MarkMovieAsWatchedAsync()
    {
        throw new NotImplementedException();
    }
    
    public Task MarkTvShowsAsWatchedAsync()
    {
        throw new NotImplementedException();
    }
    
    private async Task<ICollection<T>> GetWatchedAsync<T>(TraktWatchedType type)
    {
        using HttpClient httpClient = new();
        httpClient.BaseAddress = _config.BaseUrl;
        httpClient.DefaultRequestHeaders.Add("Accept", "application/json");
        httpClient.DefaultRequestHeaders.Add("trakt-api-version", _config.ApiVersion);
        httpClient.DefaultRequestHeaders.Add("trakt-api-key", _config.ClientId);

        try
        {
            var tmp = $"/users/lfmachadodasilva/watched/{type.ToString().ToLower()}";
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