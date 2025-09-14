package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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

	myApp := app.New()
	myWindow := myApp.NewWindow("TubeLoader")
	myWindow.Resize(fyne.NewSize(500, 500))

	// Format Selection
	formatSelect := widget.NewRadioGroup([]string{"MP3", "MP4"}, func(s string) {})
	formatSelect.SetSelected("MP3")

	// URL Entry
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("URL eingeben...")

	// Time Range Entries
	startTimeEntry := widget.NewEntry()
	startTimeEntry.SetPlaceHolder("Startzeit (z.B. 00:01:30)")
	startTimeEntry.Disable()
	endTimeEntry := widget.NewEntry()
	endTimeEntry.SetPlaceHolder("Endzeit (z.B. 00:05:45)")
	endTimeEntry.Disable()

	// Checkbox for time range
	timeRangeCheck := widget.NewCheck("Mit Time Range", func(checked bool) {
		if checked {
			startTimeEntry.Enable()
			endTimeEntry.Enable()
			return
		}
		startTimeEntry.Disable()
		endTimeEntry.Disable()
	})

	// confirm button
	confirmButton := widget.NewButton("Download starten", func() {
		format := formatSelect.Selected
		url := urlEntry.Text
		withTimeRange := timeRangeCheck.Checked
		startTime := startTimeEntry.Text
		endTime := endTimeEntry.Text

		slog.Info("Format selected", slog.String("format", format))
		slog.Info("URL", slog.String("url", url))
		slog.Info("Time Range", slog.Bool("withTimeRange", withTimeRange))
		slog.Info("Start", slog.String("startTime", startTime))
		slog.Info("End", slog.String("endTime", endTime))

		downloadConfiguration := domain.NewDownloadConfiguration(domain.FileType(format), url, withTimeRange, startTime, endTime)

		handleDownloadRequest(downloadConfiguration)
	})

	content := container.NewVBox(
		widget.NewLabel("Format:"),
		formatSelect,
		widget.NewLabel("URL:"),
		urlEntry,
		timeRangeCheck,
		container.NewGridWithColumns(2,
			container.NewVBox(
				widget.NewLabel("Startzeit:"),
				startTimeEntry,
			),
			container.NewVBox(
				widget.NewLabel("Endzeit:"),
				endTimeEntry,
			),
		),
		layout.NewSpacer(),
		confirmButton,
		layout.NewSpacer(),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
	tidyUp()

	return nil
}

// handleDownloadRequest handles the download request with the given configuration
func handleDownloadRequest(downloadConfiguration *domain.DownloadConfiguration) error {
	slog.Info("Handling download request")

	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dlCommand := ytdlp.New().
		SetWorkDir(rootPath)

	converterService, err := converter.NewService(dlCommand)
	if err != nil {
		return err
	}

	switch downloadConfiguration.Format {
	case domain.Mp3:
		err := converterService.DownloadPlaylistAsMp3()
		if err != nil {
			return err
		}
	case domain.Mp4:
		err := converterService.DownloadVideoSection(downloadConfiguration.Url, downloadConfiguration.Start, downloadConfiguration.End)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown format: %s", downloadConfiguration.Format)
	}

	return nil
}

// tidyUp does some cleanup
func tidyUp() {
	fmt.Println("Exited")
}
