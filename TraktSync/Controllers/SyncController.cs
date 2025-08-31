using Microsoft.AspNetCore.Mvc;
using TraktSync.Handler;

namespace TraktSync.Controllers;

[ApiController]
[Route("api/[controller]")]
public class SyncController(SyncHandler syncHandler) : ControllerBase
{
    [HttpPost]
    [ProducesResponseType(StatusCodes.Status200OK)]
    public async Task<IActionResult> SyncAsync(CancellationToken cancellationToken = default)
    {
        await syncHandler.SyncAsync(cancellationToken);
        return Ok();
    }
}