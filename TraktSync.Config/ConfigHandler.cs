namespace TraktSync.Config;

public class ConfigHandler
{
    public Config GetAsync()
    {
        return new Config
        {
            Trakt = new TraktConfig
            {
                UserName = "lfmachadodasilva",
                ClientId = "eb4ede9a384157e9aa60aad8c72c36c0485215659c82ad7b1fe965359a55caf4"
            },
            Emby = new EmbyConfig
            {
                BaseUrl = new Uri("http://192.168.1.13:8096"),
                ApiKey = "b039ba2b065e4ba1bca2307cce593478",
                UserId = "aac3a78d9f184ea480fb1629e76aad57"
            }
        };
    }
}