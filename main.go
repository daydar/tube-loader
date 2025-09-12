package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/lrstanley/go-ytdlp"
)

func main() {
	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)
	ytdlp.MustInstallFFmpeg(context.TODO(), nil)

	dl := ytdlp.New().
		PresetAlias("mp3").                // Use the "mp3" preset
		YesPlaylist().                     // Download the whole playlist if it is a video with a playlist parameter
		Output("output/%(title)s.%(ext)s") // Output to the "output" directory

	data, err := os.ReadFile("song_links.json")
	if err != nil {
		panic(err)
	}

	var SongLinks struct {
		Songs []string `json:"songs"`
	}

	err = json.Unmarshal(data, &SongLinks)
	if err != nil {
		panic(err)
	}

	for _, url := range SongLinks.Songs {
		_, err := dl.Run(context.TODO(), url)
		if err != nil {
			panic(err)
		}
	}
}
