namespace TraktSync.Trakt.Models;

public class TraktConfig
{
    public Uri BaseUrl { get; set; } = new("https://api.trakt.tv");
    public string ApiVersion { get; set; } = "2";
    public string RedirectUrL { get; set; } = "urn:ietf:wg:oauth:2.0:oob";
    
    public string? UserName { get; set; }
    public string? ClientId { get; set; }
    public string? ClientSecret { get; set; }
    public string? AccessToken { get; set; }
    public string? RefreshToken { get; set; }
    public string? Code { get; set; }
    public DateTime? ExpiresIn { get; set; }
    public DateTime? CreatedAt { get; set; }
}