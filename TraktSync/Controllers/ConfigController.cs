using Microsoft.AspNetCore.Mvc;
using TraktSync.Config;

namespace TraktSync.Controllers;

[ApiController]
[Route("api/[controller]")]
public class ConfigController(ConfigHandler configHandler) : ControllerBase
{
    [HttpGet]
    public IActionResult GetAsync() =>
        Ok(configHandler.GetAsync());
    
    [HttpPost]
    public IActionResult UpdateConfig([FromBody] Config.Config config) =>
        Ok(configHandler.UpdateConfig(config));
}