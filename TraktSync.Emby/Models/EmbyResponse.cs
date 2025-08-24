using System.Text.Json.Serialization;

namespace TraktSync.Emby.Models;

public class EmbyResponse
{
    [JsonPropertyName("Items")]
    public ICollection<EmbyItemResponse>? Items { get; set; }
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
    [JsonPropertyName("LastPlayedDate")]
    public DateTime? LastPlayedDate { get; set; }
}

public class EmbyItemResponse
{
    [JsonPropertyName("Id")]
    public required string Id { get; set; }
    [JsonPropertyName("Name")]
    public required string Name { get; set; }
    [JsonPropertyName("Type")]
    public string? Type { get; set; }
    [JsonPropertyName("ProviderIds")] 
    public EmbyItemIdsResponse? Ids { get; set; }
    [JsonPropertyName("UserData")]
    public EmbyItemDataResponse? Data { get; set; }
    [JsonPropertyName("SeriesId")]
    public string? ParentId { get; set; }
    [JsonPropertyName("RunTimeTicks")]
    public long RunTimeTicks { get; set; }
    [JsonPropertyName("IndexNumber")]
    public short? Episode { get; set; }
    [JsonPropertyName("ParentIndexNumber")]
    public short? Season { get; set; }

    public ICollection<EmbyItemResponse>? Episodes { get; set; }
}