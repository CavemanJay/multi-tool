using System;

namespace MultiTool.Cli
{
    public class ExeNotFoundException : Exception
    {
        public ExeNotFoundException(string exe) : base($"{exe} executable not found in path.")
        {
        }
    }
}