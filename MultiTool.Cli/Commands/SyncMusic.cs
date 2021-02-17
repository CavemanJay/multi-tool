using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using CliFx;
using CliFx.Attributes;
using MultiTool.Core;
using MultiTool.MusicSync;
using MultiTool.MusicSync.Models;
using MultiTool.Prompt;
using static MultiTool.MusicSync.Youtube;

namespace MultiTool.Cli.Commands
{
    [Command("music", Description = "Synchronize a playlist from youtube")]
    public class SyncMusicCommand : ICommand
    {
        [CommandOption("secret-file", 's',
            Description = "The PATH to the json credentials file provided by youtube api.")]
        public string SecretFilePath { get; set; }

        [CommandOption("playlist", 'p',
            Description = "The playlist to download. Will ask if not provided.")]
        public string? PlayListName { get; set; } = null;

        [CommandOption("limit", 'l', Description = "The maximum number of songs to download.")]
        public int? Limit { get; set; } = null;

        [CommandOption("root", 'r', Description = "The output path to download playlists to.")]
        public string SyncRoot { get; set; }

        [CommandOption("exe", 'e', Description = "Custom exe path for youtube-dl.")]
        public string? ExePath { get; set; } = null;

        [CommandOption("apikey", 'k',
            Description = "Youtube API Key. Not needed if the secrets file exists.")]
        public string? ApiKey { get; set; } = null;

        private YTClient _client;

        public SyncMusicCommand()
        {
            SyncRoot = Path.Join(GetHomePath(), "Sync", "Music");

            SecretFilePath =
                Path.Join(Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData),
                    "MultiTool", "client_secret.json");
        }

        private static string GetHomePath() =>
            (Environment.OSVersion.Platform == PlatformID.Unix
                ? Environment.GetEnvironmentVariable("HOME")
                : Environment.ExpandEnvironmentVariables("%HOMEDRIVE%%HOMEPATH%")) ??
            throw new Exception("Unable to locate home path.");

        private static string GetPlaylist(IEnumerable<Playlist> playlists)
        {
            var result = AutoPrompt.PromptForInput(
                "Which playlist would you like to download? (Use arrow keys to choose) ",
                playlists.Select(pl => pl.Name), false);

            return result ?? throw new Exception("Playlist name cannot be blank");
        }

        private Exception? DownloadPlaylist(Playlist playlist)
        {
            var outPath = Path.Join(SyncRoot, playlist.Name);
            Directory.CreateDirectory(outPath);
            var videos = _client.getVideos(playlist).ToList();

            // Console.WriteLine(JsonConvert.SerializeObject(videos.Select(x => x.Title).ToList()));

            try
            {
                var units = Ytdl
                    .DownloadPlaylist(new Ytdl.DownloadPlaylistOptions(outPath, videos), Limit)
                    .ToList();
                units.ForEach(Console.WriteLine);

                return null;
            }
            catch (Exception ex)
            {
                return ex;
            }
        }

        public ValueTask ExecuteAsync(IConsole console)
        {
            Process.Start("/bin/bash","-c 'env'").WaitForExit();
            // var youtubeDlPath =
            //     ExePath?.Replace("~", GetHomePath()) ?? Utils.GetExePath("youtube-dl");
            // if (youtubeDlPath is null)
            // {
            //     return ValueTask.FromException(new ExeNotFoundException("youtube-dl"));
            // }
            //
            // if (Utils.GetExePath("ffmpeg") is null)
            // {
            //     return ValueTask.FromException(new ExeNotFoundException("ffmpeg"));
            // }
            //
            // var secretsFile = new FileInfo(SecretFilePath);
            // if (!secretsFile.Exists && ApiKey is null)
            // {
            //     return ValueTask.FromException(
            //         new Exception("ApiKey cannot be empty if secrets file does not exist"));
            // }
            //
            // var credential = secretsFile.Exists
            //     ? YoutubeCredential.NewJsonFile(SecretFilePath)
            //     : YoutubeCredential.NewApiKey(ApiKey);
            //
            // _client = new YTClient(credential);
            //
            // var playlists = _client.getPlaylists() ??
            //                 throw new Exception("Unable to retrieve playlists");
            //
            // PlayListName ??= GetPlaylist(playlists);
            //
            // var chosen = playlists.Single(x => x.Name == PlayListName)!;
            //
            // // Ytdl.ExePath = youtubeDlPath;
            // var error = DownloadPlaylist(chosen);

            // return error is null ? ValueTask.CompletedTask : ValueTask.FromException(error);
            return default;
        }
    }
}