using TickerQ.Utilities.Base;
using TickerQ.Utilities.Models;
using TraktSync.Handler;

namespace TraktSync.CronJob;

public class AutoSyncJob(SyncHandler syncHandler, ILogger<AutoSyncJob> logger)
{
    [TickerFunction(functionName: "AutoSyncJob", cronExpression: "0 */12 * * *")]
    public async Task Run(TickerFunctionContext<string> tickerContext, CancellationToken cancellationToken)
    {
        try
        {
            await syncHandler.SyncAsync(cancellationToken);
        }
        catch (Exception ex)
        {
            logger.LogError(ex, "AutoSyncJob: failed to complete");
        }
    }
}