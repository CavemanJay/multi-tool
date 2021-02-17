namespace MultiTool.MusicSync.Models

type Playlist =
    { Name: string
      Id: string }
    member this.link() =
        sprintf "https://youtube.com/playlist?list=%s" this.Id

type Video =
    { Id: string
      Title: string }
    member this.link() =
        sprintf "https://youtube.com/watch?v=%s" this.Id
