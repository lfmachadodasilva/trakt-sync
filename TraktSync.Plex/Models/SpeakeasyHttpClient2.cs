using System.Net.Http.Headers;
using LukeHagar.PlexAPI.SDK.Utils;

namespace TraktSync.Plex.Models;

internal sealed class SpeakeasyHttpClient2 : SpeakeasyHttpClient, IDisposable
{
    public SpeakeasyHttpClient2()
    {
        httpClient.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/json"));
    }

    public void Dispose() => httpClient.Dispose();
}