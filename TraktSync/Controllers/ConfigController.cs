using Microsoft.AspNetCore.Mvc;
using TraktSync.Config;

namespace TraktSync.Controllers;

[ApiController]
[Route("api/[controller]")]
public class ConfigController(ConfigHandler configHandler) : ControllerBase
{
    [HttpGet]
    [ProducesResponseType(typeof(Config.Config), StatusCodes.Status200OK)]
    public IActionResult GetAsync() => Ok(configHandler.GetAsync());
    
    [HttpPost]
    [ProducesResponseType(typeof(Config.Config), StatusCodes.Status200OK)]
    public IActionResult UpdateConfig([FromBody] Config.Config config) =>
        Ok(configHandler.UpdateConfig(config));
}