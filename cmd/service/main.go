package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/lrstanley/go-ytdlp"
	"sundrop.com/tube-loader/pkg/converter"
	"sundrop.com/tube-loader/pkg/domain"
)

const (
	songLinksPath = "song_links.json"
)

func main() {
	slog.Info("Starting tube loader")

	if err := run(context.TODO()); err != nil {
		slog.Error("error while running tube loader", err.Error())
	}
}

func run(ctx context.Context) error {

	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dlCommand := ytdlp.New().
		SetWorkDir(rootPath)

	converterService, err := converter.NewService(dlCommand)
	if err != nil {
		panic(err)
	}

	url := ""
	startTime := "27:00"
	endTime := "28:00"

	err = converterService.DownloadVideoSection(url, startTime, endTime)
	if err != nil {
		panic(err)
	}

	// DownloadPlaylistAsMp3()
	return nil
}

func DownloadPlaylistAsMp3() {
	ValidateRequirements(domain.Mp3)

	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dlCommand := ytdlp.New().
		SetWorkDir(rootPath).
		PresetAlias(domain.Mp3.String()).  // Use the "mp3" preset
		YesPlaylist().                     // Download the whole playlist if it is a video with a playlist parameter
		Output("output/%(title)s.%(ext)s") // Output to the "output" directory

	songLinksPath := rootPath + "/" + songLinksPath
	data, err := os.ReadFile(songLinksPath)
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
		_, err := dlCommand.Run(context.TODO(), url)
		if err != nil {
			panic(err)
		}
	}
}

func ValidateRequirements(fileType domain.FileType) {
	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)
	if fileType == domain.Mp3 {
		// If ffmpeg isn't installed yet, download and cache it for further use.
		ytdlp.MustInstallFFmpeg(context.TODO(), nil)
	}
}
