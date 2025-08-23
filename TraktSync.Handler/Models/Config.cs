using TraktSync.Emby.Models;
using TraktSync.Plex.Models;
using TraktSync.Trakt.Models;

namespace TraktSync.Handler.Models;

public class Config
{
    public TraktConfig? TraktConfig { get; set; }
    public PlexConfig? PlexConfig { get; set; }
    public EmbyConfig? EmbyConfig { get; set; }
}