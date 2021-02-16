using System.Threading.Tasks;
using CliFx;
using CliFx.Attributes;

namespace MultiTool.Cli.Commands
{
    [Command("listen", Description = "Listen on the specified port for a client")]
    public class ListenCommand : ICommand
    {
        [CommandOption("port", 'p', Description = "The port to listen on.")]
        public int Port { get; set; } = 8081;

        public ValueTask ExecuteAsync(IConsole console)
        {
            console.Output.WriteLine($"Listening on port {Port}");

            return default;
        }
    }
}