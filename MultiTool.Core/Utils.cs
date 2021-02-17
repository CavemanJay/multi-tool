using System;
using System.IO;
using System.Linq;

namespace MultiTool.Core
{
    public class Utils
    {
        public static string? GetExePath(string exeName)
        {
            if (File.Exists(exeName))
                return Path.GetFullPath(exeName);

            var values = Environment.GetEnvironmentVariable("PATH")!;
            return values?.Split(Path.PathSeparator).Select(path => Path.Combine(path, exeName))
                .FirstOrDefault(File.Exists);
        }
    }
}