package converter

import "github.com/lrstanley/go-ytdlp"

type Service struct {
	command *ytdlp.Command
}

// NewService creates a new Service with the given command.
// The command is passed by pointer because it has a RWMutex lock field
func NewService(command *ytdlp.Command) (*Service, error) {
	return &Service{
		command: command,
	}, nil

}

func (s *Service) DownloadVideoSection() error {
	return nil
}
