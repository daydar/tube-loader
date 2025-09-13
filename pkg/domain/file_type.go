package domain

// FileType is the type of the file
type FileType string

const (
	// Mp3 is the type for files with mp3 extension
	Mp3 FileType = "mp3"
	// Mp4 is the type for files with mp4 extension
	Mp4 FileType = "mp4"
)

// String returns the string representation of the FileType
func (fileType FileType) String() string {
	return string(fileType)
}
