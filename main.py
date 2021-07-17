from pickletools import pyint
import youtube_dl
from PyInquirer import prompt

channel_id = 'UC72n1Xy0MT1Dy6c_dV-chMA'
ytdl = youtube_dl.YoutubeDL({'extract_flat': True})
with ytdl:
    playlists = ytdl.extract_info(
        F"https://www.youtube.com/channel/{channel_id}/playlists", download=False)

playlist_names = [x['title'] for x in playlists['entries']]

questions = [
    {
        'type': 'list',
        'name': 'mode',
        'message': 'Mode:',
        'choices': ['music', 'files']
    },
    {
        'type': 'list',
        'name': 'playlist',
        'message': 'Playlist:',
        'choices': playlist_names
    },
]

answers = prompt(questions)

pass
