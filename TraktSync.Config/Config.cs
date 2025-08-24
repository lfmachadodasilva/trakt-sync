namespace TraktSync.Config;

public class TraktConfig
{
    public Uri BaseUrl { get; set; } = new("https://api.trakt.tv");
    public string ApiVersion { get; set; } = "2";
    public string RedirectUrl { get; set; } = "urn:ietf:wg:oauth:2.0:oob";
    public string? UserName { get; set; }
    public string? ClientId { get; set; }
    public string? ClientSecret { get; set; }
    public string? AccessToken { get; set; }
    public string? RefreshToken { get; set; }
}

public class EmbyConfig
{
    public Uri? BaseUrl { get; set; }
    public string? UserId { get; set; }
    public string? ApiKey { get; set; }
}

public class PlexConfig
{
    // TODO: Implement Plex config
}

public class Config
{
    public TraktConfig? Trakt { get; set; }
    public EmbyConfig? Emby { get; set; }
    public PlexConfig? Plex { get; set; }
}