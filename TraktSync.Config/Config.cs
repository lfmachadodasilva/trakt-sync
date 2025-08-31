using System.ComponentModel.DataAnnotations;

namespace TraktSync.Config;

public class TraktConfig
{
    public Uri BaseUrl { get; set; } = new("https://api.trakt.tv");
    public string ApiVersion { get; set; } = "2";
    public string RedirectUrl { get; set; } = "urn:ietf:wg:oauth:2.0:oob";
    [Required]
    [MaxLength(64)]
    public string ClientId { get; set; } = "CLIENT_ID";
    [Required]
    [MaxLength(64)]
    public string ClientSecret { get; set; } = "CLIENT_SECRET";
    [MaxLength(64)]
    public string? AccessToken { get; set; }
    [MaxLength(64)]
    public string? RefreshToken { get; set; }
    [MaxLength(16)]
    public string? TokenType { get; set; }
}

public class EmbyConfig
{
    [Required]
    public Uri BaseUrl { get; set; } = new("https://localhost:8096");
    [Required]
    [MaxLength(32)]
    public string UserId { get; set; } = "USER_IDE";
    [Required]
    [MaxLength(32)]
    public string ApiKey { get; set; } = "API_KEY";
}

public class PlexConfig
{
    [Required]
    public Uri BaseUrl { get; set; } = new("https://localhost:32400");
    [Required]
    [MaxLength(64)]
    public string ApiKey { get; set; } = "API_KEY";
}

public class Config
{
    [Required]
    public TraktConfig Trakt { get; set; } = new();
    public EmbyConfig? Emby { get; set; } = new();
    public PlexConfig? Plex { get; set; } = new();
}