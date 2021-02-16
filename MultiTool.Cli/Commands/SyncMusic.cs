using System.Threading.Tasks;
using CliFx;
using CliFx.Attributes;

namespace MultiTool.Cli.Commands
{
    [Command("music", Description = "")]
    public class SyncMusicCommand : ICommand
    {
        public ValueTask ExecuteAsync(IConsole console)
        {
            return default;
        }
    }
}