package converter

import (
	"context"
	"log/slog"
	"strconv"

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
		PresetAlias(domain.Mp4.String()).  // Use the "mp4" preset
		Output("output/%(title)s.%(ext)s") // Output to the "output" directory

	return &Service{
		command: command,
	}, nil

}

func (s *Service) DownloadVideoSection(url string) error {

	s.command.DownloadSections("*27:48-27:56")

	result, err := s.command.Run(context.TODO(), url)
	if err != nil {
		return err
	}
	slog.Info("result with exit code", strconv.Itoa(result.ExitCode))

	return nil
}
