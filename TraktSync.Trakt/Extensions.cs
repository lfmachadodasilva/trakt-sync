using Microsoft.Extensions.DependencyInjection;

namespace TraktSync.Trakt;

public static class Extensions
{
    public static IServiceCollection AddTrakt(this IServiceCollection services) =>
        services.AddSingleton<TraktClient>();
}