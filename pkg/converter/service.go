package converter

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/lrstanley/go-ytdlp"
	"sundrop.com/tube-loader/pkg/domain"
)

type Service struct {
	command *ytdlp.Command
}

// NewService creates a new Service with the given command.
// The command is passed by pointer because it has a RWMutex lock field
func NewService() (*Service, error) {

	rootPath, err := os.Getwd()
	if err != nil {
		slog.Error("error while getting current working directory", slog.String("error", err.Error()))
		return nil, err
	}

	command := ytdlp.New().
		SetWorkDir(rootPath).
		Output("output/%(title)s.%(ext)s") // Output to the "output" directory

	return &Service{
		command: command,
	}, nil

}

// Download downloads the song with the given configuration
func (s *Service) Download(downloadConfiguration *domain.DownloadConfiguration) error {
	slog.Info("Starting Download...")

	validateRequirements(downloadConfiguration.Format)

	s.command.PresetAlias(downloadConfiguration.Format.String()) // Use the preset alias for the given file type

	if downloadConfiguration.WithTimeRange {
		timeRangeRegex, err := generateTimeRangeRegex(downloadConfiguration.Start, downloadConfiguration.End)
		if err != nil {
			slog.Error("error while generating time range regex", slog.String("error", err.Error()))
			return err
		}

		s.command.DownloadSections(timeRangeRegex)
	}

	if downloadConfiguration.WithPlaylist {
		s.command.YesPlaylist()
	}

	slog.Info("Downloading", slog.String("url", downloadConfiguration.Url))

	result, err := s.command.Run(context.TODO(), downloadConfiguration.Url)
	if err != nil {
		slog.Error("error while downloading", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Download finished")

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

	startEndtimeRange := regexp.MustCompile(`^(?:(?:[01]\d|2[0-3]):[0-5]\d(?::[0-5]\d)?\s*-\s*(?:[01]\d|2[0-3]):[0-5]\d(?::[0-5]\d)?)(?:\s*,\s*(?:[01]\d|2[0-3]):[0-5]\d(?::[0-5]\d)?\s*-\s*(?:[01]\d|2[0-3]):[0-5]\d(?::[0-5]\d)?)*$`)

	if !startEndtimeRange.MatchString(startTime + "-" + endTime) {
		return "", fmt.Errorf("invalid start and end time format")
	}

	regex := "*" + startTime + "-" + endTime
	return regex, nil
}
