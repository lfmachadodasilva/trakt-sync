using System.Text.Json;
using Microsoft.Extensions.Caching.Memory;

namespace TraktSync.Config;

public class ConfigHandler(IMemoryCache cache)
{
    private const string FilePath = "./data/config.json";
    
    public Config? GetAsync()
    {
        return cache.GetOrCreate<Config?>("config", entry =>
        {
            entry.AbsoluteExpirationRelativeToNow = TimeSpan.FromHours(1);
            
            if (!File.Exists(FilePath))
            {
                return UpdateConfig(new Config());
            }

            var json = File.ReadAllText(FilePath);
            var config = JsonSerializer.Deserialize<Config>(json);
        
            cache.Set("config", config, TimeSpan.FromHours(1));
        
            return config;
        });
    }

    public Config UpdateConfig(Config config)
    {
        ArgumentNullException.ThrowIfNull(config, nameof(config));
        
        var options = new JsonSerializerOptions { WriteIndented = true };
        var directory = Path.GetDirectoryName(FilePath);
        if (!string.IsNullOrEmpty(directory) && !Directory.Exists(directory))
        {
            Directory.CreateDirectory(directory);
        }
        var json = JsonSerializer.Serialize(config, options);
        File.WriteAllText(FilePath, json);
        
        cache.Set("config", config, TimeSpan.FromHours(1));
        
        return config;
    }
}