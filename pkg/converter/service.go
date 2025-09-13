package converter

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/lrstanley/go-ytdlp"
	"sundrop.com/tube-loader/pkg/domain"
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
