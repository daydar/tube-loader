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
	myWindow := myApp.NewWindow("Tube loader")
	myWindow.Resize(fyne.NewSize(500, 500))

    // Format Selection
    formatSelect := widget.NewRadioGroup([]string{"MP3", "MP4"}, func(selected string) {
		//TODO
	})
    formatSelect.SetSelected("MP3")

    // URL Entry
    urlEntry := widget.NewEntry()
    urlEntry.SetPlaceHolder("URL eingeben...")

    // Checkbox for time range
    timeRangeCheck := widget.NewCheck("Mit Time Range", func(checked bool) {
        // Hier könntest du die Time Range Felder ein-/ausblenden
    })

    // Time Range Entries
    startTimeEntry := widget.NewEntry()
    startTimeEntry.SetPlaceHolder("Startzeit (z.B. 00:01:30)")
    endTimeEntry := widget.NewEntry()
    endTimeEntry.SetPlaceHolder("Endzeit (z.B. 00:05:45)")

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
    })

    // Layout mit Grid für bessere Anordnung
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

	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dlCommand := ytdlp.New().
		SetWorkDir(rootPath)

	converterService, err := converter.NewService(dlCommand)
	if err != nil {
		panic(err)
	}

	url := ""
	startTime := "27:00"
	endTime := "28:00"

	err = converterService.DownloadVideoSection(url, startTime, endTime)
	if err != nil {
		panic(err)
	}

	// DownloadPlaylistAsMp3()
	return nil
}

func tidyUp() {
	fmt.Println("Exited")
}
