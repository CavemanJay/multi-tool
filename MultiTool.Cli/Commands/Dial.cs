using System.Threading.Tasks;
using CliFx;
using CliFx.Attributes;

namespace MultiTool.Cli.Commands
{
    [Command("dial", Description = "Connect to an existing server instance")]
    public class DialCommand : ICommand
    {
        [CommandOption("port", 'p', Description = "The port to connect to.")]
        public int Port { get; set; } = 8081;

        [CommandParameter(0, Description = "The host to connect to.")]
        public string Host { get; set; }

        public ValueTask ExecuteAsync(IConsole console)
        {
            return default;
        }
    }
}