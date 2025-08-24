using System.Text.Json.Serialization;

namespace TraktSync.Trakt.Models;

public class TraktMarkAsWatchedIdsRequest
{
    [JsonPropertyName("imdb")]
    public string? Imdb { get; set; }
}

public class TraktMarkAsWatchedMovieRequest
{
    [JsonPropertyName("ids")]
    public TraktMarkAsWatchedIdsRequest? Ids { get; set; }
    [JsonPropertyName("watched_at")]
    public DateTime? WatchedAt { get; set; }
}

public class TraktMarkAsWatchedEpisodeRequest
{
    [JsonPropertyName("number")]
    public short? Number { get; set; }
    [JsonPropertyName("watched_at")]
    public DateTime? WatchedAt { get; set; }
}

public class TraktMarkAsWatchedSeasonRequest
{
    [JsonPropertyName("number")]
    public short? Number { get; set; }
    [JsonPropertyName("episodes")]
    public ICollection<TraktMarkAsWatchedEpisodeRequest>? Episodes { get; set; }
}

public class TraktMarkAsWatchedTvShowRequest
{
    [JsonPropertyName("ids")]
    public TraktMarkAsWatchedIdsRequest? Ids { get; set; }
    [JsonPropertyName("seasons")]
    public ICollection<TraktMarkAsWatchedSeasonRequest>? Seasons { get; set; }
}

public class TraktMarkAsWatchedRequest
{
    [JsonPropertyName("movies")]
    public ICollection<TraktMarkAsWatchedMovieRequest> Movies { get; set; } = [];
    [JsonPropertyName("shows")]
    public ICollection<TraktMarkAsWatchedTvShowRequest> TvShows { get; set; } = [];
}