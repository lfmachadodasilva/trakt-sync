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
    [HttpPost("auth")]
    public async Task<IActionResult> AuthAsync([FromBody] TraktAuthCode request)
    {
        var response = await traktClient.AuthAsync(request.Code);
        return Ok(response);
    }
    
    [HttpPost("refresh-token")]
    public async Task<IActionResult> AuthRefreshTokenAsync()
    {
        var response = await traktClient.AuthRefreshAccessTokenAsync();
        return Ok(response);
    }
}