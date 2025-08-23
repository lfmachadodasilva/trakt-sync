using System.Text.Json.Serialization;

namespace TraktSync.Trakt.Models;

public enum TraktWatchedType
{
    Movies,
    Shows
}

public class TraktWatchedIdsResponse
{
    [JsonPropertyName("trakt")]
    public int? Trakt { get; set; }
    [JsonPropertyName("slug")]
    public string? Slug { get; set; }
    [JsonPropertyName("imdb")]
    public string? Imdb { get; set; }
}

public class TraktWatchedItemResponse
{
    [JsonPropertyName("title")]
    public string? Title { get; set; }
    [JsonPropertyName("year")]
    public short? Year { get; set; }
    [JsonPropertyName("ids")]
    public TraktWatchedIdsResponse? Ids { get; set; }
}

public class TraktWatchedSeasonResponse
{
    [JsonPropertyName("number")]
    public short? Number { get; set; }
    [JsonPropertyName("last_watched_at")]
    public DateTime? LastWatchedAt { get; set; }
}

public abstract class TraktWatchedBaseResponse
{
    [JsonPropertyName("last_watched_at")]
    public DateTime? LastWatchedAt { get; set; }
    [JsonPropertyName("last_update_at")]
    public DateTime? LastUpdatedAt { get; set; }
}

public class TraktWatchedTvShowResponse : TraktWatchedBaseResponse
{
    [JsonPropertyName("show")]
    public TraktWatchedItemResponse? Show { get; set; }
    [JsonPropertyName("seasons")]
    public ICollection<TraktWatchedSeasonResponse>? Seasons { get; set; }
}

public class TraktWatchedMoviesResponse : TraktWatchedBaseResponse
{
    [JsonPropertyName("movie")]
    public TraktWatchedItemResponse? Movie { get; set; }
}