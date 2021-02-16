using System.Collections.Generic;
using System.Threading.Tasks;
using CliFx;

namespace MultiTool.Cli
{
    public class MultiTool
    {
        static async Task<int> Main() =>
            await new CliApplicationBuilder().AddCommandsFromThisAssembly()
                .UseExecutableName("multi").Build().RunAsync();
    }
}