using Microsoft.Extensions.DependencyInjection;

namespace TraktSync.Plex;

public static class Extensions
{
    public static IServiceCollection AddPlex(this IServiceCollection services) =>
        services.AddSingleton<IPlexClient, PlexClient>();
    
    public static string GetImdb(this ICollection<string> ids) =>
        ids?.FirstOrDefault(x => x.StartsWith("imdb://"))?.Replace("imdb://", string.Empty) ?? string.Empty;
    public static string GetImdb(this IEnumerable<string> ids) => ids?.ToList()?.GetImdb() ?? string.Empty;
}