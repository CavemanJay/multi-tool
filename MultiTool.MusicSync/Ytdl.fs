module MultiTool.MusicSync.Ytdl

open System.Diagnostics
open System.IO
open System.Reflection
open MultiTool.MusicSync.Models

type DownloadPlaylistOptions =
    { outputRoot: string
      videos: Video seq }

let private template = "%s(title)s.%(format)s"

let ExePath =
    Path.Join
        (Path.GetDirectoryName(Assembly.GetEntryAssembly().Location),
         "Binaries",
         "music")

let private getFileName (video: Video) =
    let fileName =
        video
            .Title
            .Replace("/", "_")
            .Replace("**OUT ON SPOTIFY**", "_OUT ON SPOTIFY")
            .Replace("|", "_")
            .Replace("*", "_")

    sprintf "%s.mp3" fileName

let private downloadVideo (video: Video) (outputRoot: string) =
    let args =
        [| Path.Join(outputRoot, template)
           video.link () |]

    Process.Start(ExePath, args).WaitForExit()


let private notDownloaded video root =
    let fileName = getFileName video
    let path = Path.Join(fileName, root)

    FileInfo(path).Exists |> not

let private downloadVideos options =
    options.videos
    |> Seq.map (fun vid -> downloadVideo vid options.outputRoot)
// for video in options.videos do
//     let result = downloadVideo video options.outputRoot
//     ()
//
// ()

let private downloadPlaylistWithLimit (options: DownloadPlaylistOptions) limit =
    options.videos
    |> Seq.filter (fun v -> notDownloaded v options.outputRoot)
    |> Seq.take limit
    |> List.ofSeq
    |> fun vids ->
        downloadVideos
            { videos = vids
              outputRoot = options.outputRoot }


let private downloadPlaylistNoLimit (options: DownloadPlaylistOptions) =
    options.videos
    |> Seq.filter (fun v -> notDownloaded v options.outputRoot)
    |> List.ofSeq
    |> fun vids ->
        downloadVideos
            { videos = vids
              outputRoot = options.outputRoot }


let DownloadPlaylist (options: DownloadPlaylistOptions)
                     (limit: System.Nullable<int>)
                     =


    match limit.HasValue with
    | true ->
        match limit.Value >= (options.videos |> Seq.length) with
        | true -> downloadPlaylistNoLimit options
        | false -> downloadPlaylistWithLimit options limit.Value
    | false -> downloadPlaylistNoLimit options
