using System.Text.Json.Serialization;

namespace TraktSync.Emby.Models;

public class EmbyResponse
{
    [JsonPropertyName("Items")]
    public ICollection<EmbyItemResponse> Items { get; set; }
    [JsonPropertyName("TotalRecordCount")] 
    public int TotalRecordCount { get; set; }
}

public class EmbyItemIdsResponse
{
    [JsonPropertyName("Imdb")] 
    public string? Imdb { get; set; }
}

public class EmbyItemDataResponse
{
    [JsonPropertyName("Played")]
    public bool Played { get; set; }
}

public class EmbyItemResponse
{
    [JsonPropertyName("Id")] 
    public required string Id { get; set; }
    [JsonPropertyName("Name")] 
    public required string Name { get; set; }
    [JsonPropertyName("Type")] 
    public required string Type { get; set; }
    [JsonPropertyName("ProviderIds")] 
    public EmbyItemIdsResponse? Ids { get; set; }
    [JsonPropertyName("UserData")] 
    public EmbyItemDataResponse? Data { get; set; }
}