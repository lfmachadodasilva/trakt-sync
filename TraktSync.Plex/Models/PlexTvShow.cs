namespace TraktSync.Plex.Models;

public class PlexTvShow
{
    public string? Id { get; set; }
    public string? Type { get; set; }
    public string? Imdb { get; set; }
    public string? Name { get; set; }
    
    public int? Season { get; set; }
    public int? Episode { get; set; }
    public bool? Played { get; set; }
    public DateTime? PlayedAt { get; set; }
    
    public object? Object { get; set; }
    
    public ICollection<PlexTvShow> Children { get; set; } = new List<PlexTvShow>();
}

