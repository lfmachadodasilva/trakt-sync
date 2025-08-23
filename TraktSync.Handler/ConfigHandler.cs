using TraktSync.Emby.Models;
using TraktSync.Handler.Models;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler;

public class ConfigHandler
{
    public Config GetConfigAsync()
    {
        return new Config
        {
            TraktConfig = new TraktConfig
            {
                BaseUrl = new Uri("https://api.trakt.tv"),
                UserName = "lfmachadodasilva",
                ClientId = "eb4ede9a384157e9aa60aad8c72c36c0485215659c82ad7b1fe965359a55caf4"
            },
            EmbyConfig = new EmbyConfig
            {
                BaseUrl = new Uri("http://localhost:8096"),
                ApiKey = "b039ba2b065e4ba1bca2307cce593478",
                UserId = "aac3a78d9f184ea480fb1629e76aad57"
            }
        };
    }
}