using Microsoft.AspNetCore.Mvc;
using TraktSync.Config;
using TraktSync.Trakt;

namespace TraktSync.Controllers;

public class TraktAuthCode
{
    public required string Code { get; set; }
}

[ApiController]
[Route("api/[controller]")]
public class TraktController(
    TraktClient traktClient, ConfigHandler configHandler) : ControllerBase
{
    [HttpGet("code")]
    public async Task<IActionResult> CodeAsync()
    {
        var response = await traktClient.CodeAsync();
        return Ok(response);
    }
    
    [HttpPost("auth")]
    public async Task<IActionResult> AuthAsync([FromBody] TraktAuthCode request)
    {
        var response = await traktClient.AuthAsync(request.Code);
        
        var config = configHandler.GetAsync();
        if (config.Trakt != null)
        {
            config.Trakt.AccessToken = response.AccessToken;
            config.Trakt.RefreshToken = response.RefreshToken;
            configHandler.UpdateConfig(config);
        }
        else
        {
            throw new NullReferenceException("Trakt config is null");
        }
        
        return Ok();
    }
    
    [HttpPost("refresh-token")]
    public async Task<IActionResult> AuthRefreshTokenAsync()
    {
        var response = await traktClient.AuthRefreshAccessTokenAsync();
        
        var config = configHandler.GetAsync();
        if (config.Trakt != null)
        {
            config.Trakt.AccessToken = response.AccessToken;
            config.Trakt.RefreshToken = response.RefreshToken;
            configHandler.UpdateConfig(config);
        }
        else
        {
            throw new NullReferenceException("Trakt config is null");
        }
        
        return Ok();
    }
}