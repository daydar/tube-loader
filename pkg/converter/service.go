package converter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/lrstanley/go-ytdlp"
	"sundrop.com/tube-loader/pkg/domain"
)

const (
	songLinksPath = "song_links.json"
)

type Service struct {
	command *ytdlp.Command
}

// NewService creates a new Service with the given command.
// The command is passed by pointer because it has a RWMutex lock field
func NewService(command *ytdlp.Command) (*Service, error) {
	command = command.
		Output("output/%(title)s.%(ext)s") // Output to the "output" directory

	return &Service{
		command: command,
	}, nil

}

// DownloadVideoSection downloads a video section from a given url with the given start and end time
func (s *Service) DownloadVideoSection(url string, startTime string, endTime string) error {
	slog.Info("Starting DownloadVideoSection")

	if startTime == "" || endTime == "" {
		return fmt.Errorf("startTime and endTime cannot be empty")
	}

	startEndtimeRange := regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d\s*-\s*([01]\d|2[0-3]):[0-5]\d(?:\s*,\s*([01]\d|2[0-3]):[0-5]\d\s*-\s*([01]\d|2[0-3]):[0-5]\d)*$`)

	if !startEndtimeRange.MatchString(startTime + "-" + endTime) {
		return fmt.Errorf("invalid start and end time format")
	}

	regex := "*" + startTime + "-" + endTime

	s.command.
		PresetAlias(domain.Mp4.String()). // Use the "mp4" preset
		DownloadSections(regex)
	result, err := s.command.Run(context.TODO(), url)
	if err != nil {
		return err
	}

	slog.Info("result", result.String(), "")
	return nil
}

// DownloadPlaylistAsMp3 downloads a playlist as mp3
func (s *Service) DownloadPlaylistAsMp3() error {
	ValidateRequirements(domain.Mp3)

	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	s.command.
		PresetAlias(domain.Mp3.String()). // Use the "mp3" preset
		YesPlaylist()                     // Download the whole playlist if it is a video with a playlist parameter

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
		_, err := s.command.Run(context.TODO(), url)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func ValidateRequirements(fileType domain.FileType) {
	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)
	if fileType == domain.Mp3 || fileType == domain.Mp4 {
		// If ffmpeg isn't installed yet, download and cache it for further use.
		ytdlp.MustInstallFFmpeg(context.TODO(), nil)
	}
}
