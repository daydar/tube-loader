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
		slog.Error("error while running tube loader", slog.String("error", err.Error()))
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

func tidyUp() {
	fmt.Println("Exited")
}
