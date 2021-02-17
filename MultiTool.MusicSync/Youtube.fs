module MultiTool.MusicSync.Youtube

open System.IO
open System.Reflection
open System.Threading
open Google.Apis.Auth.OAuth2
open Google.Apis.Services
open Google.Apis.Util.Store
open Google.Apis.YouTube.v3
open Google.Apis.YouTube.v3.Data
open MultiTool.MusicSync.Models

let private getRepeatable part =
    Google.Apis.Util.Repeatable<string> [| part |]

type YoutubeCredential =
    | JsonFile of string
    | ApiKey of string
    | UserCreds of UserCredential


let private scopes =
    [| YouTubeService.Scope.YoutubeReadonly |]

let private getYtInitializer (credential: YoutubeCredential) =
    let init = BaseClientService.Initializer()

    match credential with
    | UserCreds creds -> init.HttpClientInitializer <- creds
    | ApiKey key -> init.ApiKey <- key
    | _ -> failwith "Invalid credential type"

    init.ApplicationName <- Assembly.GetExecutingAssembly().GetName().Name

    init

let private getClientFromFile path =
    use stream =
        new FileStream(path, FileMode.Open, FileAccess.Read)

    let credential =
        GoogleWebAuthorizationBroker.AuthorizeAsync
            (GoogleClientSecrets.Load(stream).Secrets,
             scopes,
             "user",
             CancellationToken.None,

             FileDataStore("multi"),
             PromptCodeReceiver())
        |> Async.AwaitTask
        |> Async.RunSynchronously

    new YouTubeService(getYtInitializer (UserCreds credential))

let private getClient credential =
    match credential with
    | JsonFile path -> getClientFromFile path
    | ApiKey key -> new YouTubeService(getYtInitializer (ApiKey key))
    | _ -> failwith "Invalid credential type"

type YTClient(creds: YoutubeCredential) =
    let client = getClient creds

    member this.getPlaylists() =
        let parts = getRepeatable "snippet"

        let request = client.Playlists.List(parts)
        request.MaxResults <- 100L
        request.Mine <- true

        request.Execute().Items
        |> Seq.map (fun pl -> { Id = pl.Id; Name = pl.Snippet.Title })

    member this.getVideos(pl: Playlist) =
        let convertPlaylistItem (pli: PlaylistItem) =
            { Id = pli.Id
              Title = pli.Snippet.Title }

        let parts = getRepeatable "snippet,contentDetails"

        let request = client.PlaylistItems.List(parts)
        request.MaxResults <- 100L
        request.PlaylistId <- pl.Id

        seq {

            let mutable response = request.Execute()

            yield response.Items |> Seq.map convertPlaylistItem

            // While the next page token is not null
            while not <| isNull response.NextPageToken do
                request.PageToken <- response.NextPageToken
                response <- request.Execute()

                yield response.Items |> Seq.map convertPlaylistItem
        }
        |> Seq.concat
        |> Seq.distinct
