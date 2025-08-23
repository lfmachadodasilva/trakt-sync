using Microsoft.Extensions.DependencyInjection;
using TraktSync.Emby;
using TraktSync.Plex;
using TraktSync.Trakt;

namespace TraktSync.Handler;

public static class Extensions
{
    public static IServiceCollection AddHandler(this IServiceCollection services) =>
        services
            .AddSingleton<SyncHandler>()
            .AddSingleton<SyncTvShowsHandler>()
            .AddSingleton<SyncMoviesHandler>()
            .AddSingleton<ConfigHandler>()
            .AddTrakt()
            .AddEmby()
            .AddPlex();
}