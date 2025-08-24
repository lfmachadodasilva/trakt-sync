using System.Text.Json.Serialization;

namespace TraktSync.Trakt.Models;

public class TraktAuthCode
{
    [JsonPropertyName("client_id")]
    public string? ClientId { get; set; }
    [JsonPropertyName("redirect_uri")]
    public string? RedirectUrl { get; set; }
    [JsonPropertyName("response_type")]
    public string ResponseType { get; set; } = "code";
}

public class TraktAuthRequest
{
    [JsonPropertyName("code")]
    public string? Code { get; set; }
    [JsonPropertyName("client_id")]
    public string? ClientId { get; set; }
    [JsonPropertyName("client_secret")]
    public string? ClientSecret { get; set; }
    [JsonPropertyName("redirect_uri")]
    public string? RedirectUrl { get; set; }
    [JsonPropertyName("grant_type")]
    public required string GrantType { get; set; }
}

public class TraktAuthRefreshRequest
{
    [JsonPropertyName("refresh_token")]
    public string? RefreshToken { get; set; }
    [JsonPropertyName("client_id")]
    public string? ClientId { get; set; }
    [JsonPropertyName("client_secret")]
    public string? ClientSecret { get; set; }
    [JsonPropertyName("redirect_uri")]
    public string? RedirectUrl { get; set; }
    [JsonPropertyName("grant_type")]
    public required string GrantType { get; set; }
}

public class TraktAuthResponse
{
    [JsonPropertyName("access_token")]
    public string? AccessToken { get; set; }
    [JsonPropertyName("token_type")]
    public string? TokenType { get; set; }
    [JsonPropertyName("refresh_token")]
    public string? RefreshToken { get; set; }
    [JsonPropertyName("expires_in")]
    public long? ExpiresIn { get; set; }
    [JsonPropertyName("created_at")]
    public long? CreatedAt { get; set; }
    [JsonPropertyName("scope")]
    public string? Scope { get; set; }
}