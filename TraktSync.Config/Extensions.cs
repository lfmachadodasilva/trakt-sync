using Microsoft.Extensions.DependencyInjection;

namespace TraktSync.Config;

public static class Extensions
{
    public static IServiceCollection AddConfig(this IServiceCollection services) =>
        services.AddSingleton<ConfigHandler>();
}