using System.Threading.Tasks;
using CliFx;
using CliFx.Attributes;

namespace MultiTool.Cli.Commands
{
    [Command("archive")]
    public class ArchiveCommand : ICommand
    {
        public ValueTask ExecuteAsync(IConsole console)
        {
            return default;
        }
    }
}