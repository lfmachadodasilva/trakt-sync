using Microsoft.Extensions.DependencyInjection;

namespace TraktSync.Emby;

public static class Extensions
{
    public static IServiceCollection AddEmby(this IServiceCollection services) =>
        services.AddSingleton<IEmbyClient, EmbyClient>();
}