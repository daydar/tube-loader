package main

import (
	"context"
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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

// run runs the application
func run(ctx context.Context) error {
	setupUI()
	return nil
}

func setupUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TubeLoader")
	myWindow.Resize(fyne.NewSize(500, 600))

	// Format Selection
	formatSelect := widget.NewRadioGroup([]string{"mp3", "mp4"}, func(s string) {})
	formatSelect.SetSelected("mp3")

	// URL Entry
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter url...")

	// Time Range Entries
	startTimeEntry := widget.NewEntry()
	startTimeEntry.SetPlaceHolder("Start Time (e.g. 00:01:30)")
	startTimeEntry.Disable()
	endTimeEntry := widget.NewEntry()
	endTimeEntry.SetPlaceHolder("End Time (e.g. 00:05:45)")
	endTimeEntry.Disable()

	// Checkbox for time range
	timeRangeCheck := widget.NewCheck("With Time Range", func(checked bool) {
		if checked {
			startTimeEntry.Enable()
			endTimeEntry.Enable()
			return
		}
		startTimeEntry.Disable()
		endTimeEntry.Disable()
	})

	// Checkbox for playlist
	playlistCheck := widget.NewCheck("With Playlist", func(checked bool) {})

	// progress bar
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// info label
	resultInfoLabel := widget.NewLabel("Video downloaded at: Downloads/TubeLoader/output")
	resultInfoLabel.Hide()

	// confirm button
	confirmButton := widget.NewButton("Start Download", func() {
		format := formatSelect.Selected
		url := urlEntry.Text
		withTimeRange := timeRangeCheck.Checked
		startTime := startTimeEntry.Text
		endTime := endTimeEntry.Text
		withPlaylist := playlistCheck.Checked

		slog.Info("Format selected", slog.String("format", format))
		slog.Info("URL", slog.String("url", url))
		slog.Info("Time Range", slog.Bool("withTimeRange", withTimeRange))
		slog.Info("Start", slog.String("startTime", startTime))
		slog.Info("End", slog.String("endTime", endTime))

		downloadConfiguration := domain.NewDownloadConfiguration(domain.FileType(format), url, withTimeRange, startTime, endTime, withPlaylist)

		go func() {
			fyne.Do(func() {
				progressBar.Show()
				progressBar.SetValue(0.2)
			})
			handleDownloadRequest(downloadConfiguration)
			fyne.Do(func() {
				progressBar.SetValue(0.8)
				progressBar.SetValue(1.0)
				resultInfoLabel.Show()
			})
		}()
	})

	content := container.NewVBox(
		widget.NewLabel("Format:"),
		formatSelect,
		widget.NewLabel("URL:"),
		urlEntry,
		layout.NewSpacer(),
		timeRangeCheck,
		container.NewGridWithColumns(2,
			container.NewVBox(
				widget.NewLabel("Start Time:"),
				startTimeEntry,
			),
			container.NewVBox(
				widget.NewLabel("End Time:"),
				endTimeEntry,
			),
		),
		layout.NewSpacer(),
		playlistCheck,
		layout.NewSpacer(),
		progressBar,
		resultInfoLabel,
		layout.NewSpacer(),
		confirmButton,
		layout.NewSpacer(),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
	tidyUp()
}

// handleDownloadRequest handles the download request with the given configuration
func handleDownloadRequest(downloadConfiguration *domain.DownloadConfiguration) error {
	slog.Info("Handling download request")

	converterService, err := converter.NewService()
	if err != nil {
		slog.Error("error while creating converter service", slog.String("error", err.Error()))
		return err
	}

	err = converterService.Download(downloadConfiguration)
	if err != nil {
		slog.Error("error while downloading", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// tidyUp does some cleanup
func tidyUp() {
	fmt.Println("Exited")
}
