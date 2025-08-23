using Microsoft.Extensions.DependencyInjection;

namespace TraktSync.Plex;

public static class Extensions
{
    public static IServiceCollection AddPlex(this IServiceCollection services) =>
        services.AddSingleton<PlexClient>();
}