package converter

import (
	"context"
	"fmt"
	"log/slog"
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
func (s *Service) DownloadVideo(url string, startTime string, endTime string) error {
	slog.Info("Starting DownloadVideoSection")

	validateRequirements(domain.Mp4)

	timeRangeRegex, err := generateTimeRangeRegex(startTime, endTime)
	if err != nil {
		return err
	}
	s.command.
		PresetAlias(domain.Mp4.String()). // Use the "mp4" preset
		DownloadSections(timeRangeRegex).
		YesPlaylist()

	result, err := s.command.Run(context.TODO(), url)
	if err != nil {
		return err
	}

	slog.Info("result", result.String(), "")
	return nil
}

// DownloadAudio downloads an audio section from a given url with the given start and end time
func (s *Service) DownloadAudio(url string, startTime string, endTime string) error {
	slog.Info("Starting DownloadAudio")

	validateRequirements(domain.Mp3)

	generateTimeRangeRegex(startTime, endTime)
	timeRangeRegex, err := generateTimeRangeRegex(startTime, endTime)
	if err != nil {
		return err
	}

	s.command.
		PresetAlias(domain.Mp3.String()). // Use the "mp3" preset
		DownloadSections(timeRangeRegex).
		YesPlaylist()

	result, err := s.command.Run(context.TODO(), url)
	if err != nil {
		return err
	}

	slog.Info("result", result.String(), "")
	return nil
}

// validateRequirements checks if yt-dlp and ffmpeg are installed
func validateRequirements(fileType domain.FileType) {
	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)
	if fileType == domain.Mp3 || fileType == domain.Mp4 {
		// If ffmpeg isn't installed yet, download and cache it for further use.
		ytdlp.MustInstallFFmpeg(context.TODO(), nil)
	}
}

// generateTimeRangeRegex generates a regex for the given start and end time
func generateTimeRangeRegex(startTime string, endTime string) (string, error) {
	if startTime == "" || endTime == "" {
		return "", fmt.Errorf("startTime and endTime cannot be empty")
	}

	startEndtimeRange := regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d\s*-\s*([01]\d|2[0-3]):[0-5]\d(?:\s*,\s*([01]\d|2[0-3]):[0-5]\d\s*-\s*([01]\d|2[0-3]):[0-5]\d)*$`)

	if !startEndtimeRange.MatchString(startTime + "-" + endTime) {
		return "", fmt.Errorf("invalid start and end time format")
	}

	regex := "*" + startTime + "-" + endTime
	return regex, nil
}
