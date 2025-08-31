using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Mvc;
using TraktSync.Trakt;

namespace TraktSync.Controllers;

public class TraktAuthCode
{
    [Required]
    [MaxLength(8)]
    public required string? Code { get; set; }
}

[ApiController]
[Route("api/[controller]")]
public class TraktController(TraktClient traktClient) : ControllerBase
{
    [HttpGet("code")]
    [ProducesResponseType(typeof(string), StatusCodes.Status200OK)]
    public async Task<IActionResult> CodeAsync()
    {
        var response = traktClient.GetCodeUrl();
        await Task.CompletedTask;
        return Ok(response);
    }
    
    [HttpPost("auth")]
    [ProducesResponseType(StatusCodes.Status200OK)]
    public async Task<IActionResult> AuthAsync(
        [FromBody, Required] TraktAuthCode request,
        CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        ArgumentNullException.ThrowIfNull(request.Code);
        
        await traktClient.AuthAsync(request.Code, cancellationToken);
        return Ok();
    }
    
    [HttpPost("refresh-token")]
    [ProducesResponseType(StatusCodes.Status200OK)]
    public async Task<IActionResult> AuthRefreshTokenAsync(CancellationToken cancellationToken = default)
    {
        await traktClient.AuthRefreshAccessTokenAsync(cancellationToken);
        return Ok();
    }
}