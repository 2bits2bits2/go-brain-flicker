package main

import (
	"bytes"
	"embed"
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"log"
	"os"
	"sync/atomic"
	"time"
)

type UI struct {
	tickerDone           chan bool
	t                    *time.Ticker
	isTickerRunning      atomic.Bool
	img1                 IMG
	img2                 IMG
	mainImage            int
	flipRate             atomic.Int32
	cleanedTicker        chan bool
	rateEditor           widget.Editor
	aboutDialog          *AboutDialog
	scheduleEditor       widget.Editor
	useSchedule          bool
	schedule             []ScheduleItem
	currentScheduleIndex int
	scheduleStartTime    time.Time
}

//go:embed assets/*
var assets embed.FS

func main() {

	ui := &UI{
		tickerDone:    make(chan bool, 1),
		cleanedTicker: make(chan bool),
		// Initialize the editor with number-only filter
		rateEditor: widget.Editor{
			SingleLine: true,
			Filter:     "0123456789", // we only want numbers and at most of length of 3 as rate
			MaxLen:     3,
		},
		scheduleEditor: widget.Editor{
			SingleLine: true,
			MaxLen:     100,
		},
		aboutDialog: NewAboutDialog(),
		useSchedule: false,
		schedule:    []ScheduleItem{},
	}
	ui.flipRate.Store(1)
	ui.isTickerRunning.Store(false)

	// Load saved schedule if exists
	loadSchedule(ui)

	// Load embedded images
	var err error
	ui.img1, err = loadEmbeddedImage("assets/img1.png")
	if err != nil {
		log.Fatal(err)
	}
	ui.img2, err = loadEmbeddedImage("assets/img2.png")
	if err != nil {
		log.Fatal(err)
	}
	ui.mainImage = 1

	go func() {
		w := new(app.Window)
		w.Option(app.Title("Brain flicker"))
		w.Option(app.Size(unit.Dp(800), unit.Dp(600)))

		if err := draw(w, ui); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)

	}()
	app.Main()
}

type IMG struct {
	imgOp   paint.ImageOp
	imgSize image.Point
}

type ScheduleItem struct {
	Duration       int // in seconds
	FlickeringRate int // flips per second
	BlankTime      int // in seconds, 0 if not a blank period
}

func draw(w *app.Window, ui *UI) error {
	var ops op.Ops
	startButton := new(widget.Clickable)
	stopButton := new(widget.Clickable)
	setButton := new(widget.Clickable)
	aboutButton := new(widget.Clickable)
	scheduleButton := new(widget.Clickable)     // Add schedule button
	saveScheduleButton := new(widget.Clickable) // Add save schedule button
	useScheduleButton := new(widget.Clickable)  // Add use schedule button
	th := material.NewTheme()

	for {
		evt := w.Event()
		switch e := evt.(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if startButton.Clicked(gtx) {
				startTicker(ui, w)
			}
			if stopButton.Clicked(gtx) {
				stopTicker(ui)
			}
			if setButton.Clicked(gtx) {
				changeRate(ui, w)
			}
			if aboutButton.Clicked(gtx) {
				ui.aboutDialog.isOpen = true
			}
			if ui.aboutDialog.closeButton.Clicked(gtx) {
				ui.aboutDialog.isOpen = false
			}
			if saveScheduleButton.Clicked(gtx) {
				// Parse and save the schedule
				parseSchedule(ui, ui.scheduleEditor.Text())
				saveSchedule(ui)
			}
			if useScheduleButton.Clicked(gtx) {
				// Toggle the use schedule flag
				ui.useSchedule = !ui.useSchedule

				// If we're running, restart with the new setting
				if ui.isTickerRunning.Load() {
					stopTicker(ui)
					<-ui.cleanedTicker
					startTicker(ui, w)
				}
			}

			// Create a flex layout for the entire window
			createLayout(gtx, th, startButton, stopButton, setButton, aboutButton,
				scheduleButton, saveScheduleButton, useScheduleButton, ui)
			ui.aboutDialog.Layout(gtx, th)

			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}

func loadEmbeddedImage(path string) (IMG, error) {
	// Read file from embedded filesystem
	imgBytes, err := assets.ReadFile(path)
	if err != nil {
		return IMG{}, err
	}

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return IMG{}, err
	}

	// Convert to RGBA if it's not already

	return IMG{
		imgOp:   paint.NewImageOp(img),
		imgSize: img.Bounds().Size(),
	}, nil
}

// Load schedule from file
func loadSchedule(ui *UI) {
	// Try to read the schedule file
	data, err := os.ReadFile("schedule.txt")
	if err != nil {
		// File doesn't exist or error reading, just return without loading
		return
	}

	// Set the schedule editor text
	scheduleText := string(data)
	ui.scheduleEditor.SetText(scheduleText)

	// Parse the schedule
	parseSchedule(ui, scheduleText)
}
