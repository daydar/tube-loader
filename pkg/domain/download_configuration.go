package domain

// DownloadConfiguration is the configuration for a download
type DownloadConfiguration struct {
	Format        FileType
	Url           string
	WithTimeRange bool
	Start         string
	End           string
	WithPlaylist  bool
}

// NewDownloadConfiguration creates a new DownloadConfiguration
func NewDownloadConfiguration(format FileType, url string, withTimeRange bool, start string, end string, withPlaylist bool) *DownloadConfiguration {
	return &DownloadConfiguration{
		Format:        format,
		Url:           url,
		WithTimeRange: withTimeRange,
		Start:         start,
		End:           end,
		WithPlaylist:  withPlaylist,
	}
}
