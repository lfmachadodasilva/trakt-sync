using System.Text.Json.Serialization;

namespace TraktSync.Emby.Models;

public class EmbyRequest
{
    [JsonPropertyName("IncludeItemTypes")] 
    public string? IncludeItemTypes { get; set; }
    [JsonPropertyName("Recursive")]
    public bool Recursive { get; set; } = true;
    [JsonPropertyName("Fields")]
    public string Fields { get; set; } = "ProviderIds";
}

public enum EmbyItemType
{
    Movie,
    Series,
    Episode
}