using System.ComponentModel.DataAnnotations;

namespace TraktSync.Config;

public class TraktConfig
{
    public Uri BaseUrl { get; set; } = new("https://api.trakt.tv");
    public string ApiVersion { get; set; } = "2";
    public string RedirectUrl { get; set; } = "urn:ietf:wg:oauth:2.0:oob";
    [Required]
    [MaxLength(64)]
    public required string ClientId { get; set; }
    [Required]
    [MaxLength(64)]
    public required string ClientSecret { get; set; }
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
    [MaxLength(1024)]
    public required Uri BaseUrl { get; set; }
    [Required]
    [MaxLength(32)]
    public required string UserId { get; set; }
    [Required]
    [MaxLength(32)]
    public required string ApiKey { get; set; }
}

public class PlexConfig
{
    [Required] [MaxLength(1024)]
    public required Uri BaseUrl { get; set; }
    [Required]
    [MaxLength(64)]
    public required string ApiKey { get; set; }
}

public class Config
{
    public TraktConfig? Trakt { get; set; }
    public EmbyConfig? Emby { get; set; }
    public PlexConfig? Plex { get; set; }
}