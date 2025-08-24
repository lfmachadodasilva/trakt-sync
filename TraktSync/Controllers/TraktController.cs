using Microsoft.AspNetCore.Mvc;
using TraktSync.Trakt;

namespace TraktSync.Controllers;

public class TraktAuthCode
{
    public required string Code { get; set; }
}

[ApiController]
[Route("api/[controller]")]
public class TraktController(TraktClient traktClient) : ControllerBase
{
    [HttpGet("code")]
    public async Task<IActionResult> CodeAsync()
    {
        var response = traktClient.GetCodeUrl();
        await Task.CompletedTask;
        return Ok(response);
    }
    
    [HttpPost("auth")]
    public async Task<IActionResult> AuthAsync([FromBody] TraktAuthCode request)
    {
        await traktClient.AuthAsync(request.Code);
        return Ok();
    }
    
    [HttpPost("refresh-token")]
    public async Task<IActionResult> AuthRefreshTokenAsync()
    {
        await traktClient.AuthRefreshAccessTokenAsync();
        return Ok();
    }
}