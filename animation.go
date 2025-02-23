package main

import (
	"fmt"
	"gioui.org/app"

	"time"
)

func (ui *UI) FlipImage(w *app.Window) {
	if 1 == ui.mainImage {
		ui.mainImage = 2
	} else {
		ui.mainImage = 1
	}
	w.Invalidate()
}

func changeRate(ui *UI, w *app.Window) {
	text := ui.rateEditor.Text()
	var newRate int
	if text != "" {
		_, err := fmt.Sscanf(text, "%d", &newRate)
		if err != nil {
			newRate = 1
		}
	}
	if newRate <= 0 || newRate >= 100 {
		newRate = 1
	}
	ui.flipRate.Store(int32(newRate))
	if ui.isTickerRunning.Load() {
		stopTicker(ui)
		<-ui.cleanedTicker
		startTicker(ui, w)
	}
}

func calculateRate(flipsPerSecond int) (int, error) {
	if flipsPerSecond <= 0 {
		return 0, fmt.Errorf("flips per second must be positive, got: %d", flipsPerSecond)
	}
	if flipsPerSecond > 100 {
		return 0, fmt.Errorf("flips per second must be <= 100, got: %d", flipsPerSecond)
	}

	// Convert flips per second to milliseconds between flips
	// 1000ms / flips = ms per flip
	return 1000 / flipsPerSecond, nil
}

func startTicker(ui *UI, w *app.Window) {
	if ui.isTickerRunning.Load() == true {
		return
	}
	rate, err := calculateRate(int(ui.flipRate.Load()))
	if err != nil {
		rate = 1
	}
	if ui.t == nil {
		ui.t = time.NewTicker(time.Millisecond * time.Duration(rate))
	}

	if ui.tickerDone == nil {
		ui.tickerDone = make(chan bool, 1)

	}
	ui.isTickerRunning.Store(true)

	go func() {
		defer func() {
			close(ui.tickerDone)
			ui.t.Stop()
			ui.t = nil
			ui.tickerDone = nil
			ui.isTickerRunning.Store(false)
			ui.cleanedTicker <- true
		}()
		for {
			select {
			case <-ui.t.C:
				ui.FlipImage(w)

			case <-ui.tickerDone:
				return
			}
		}
	}()
}

// okay so we have error when
// clicked start, set, stop
// and then start and set

func stopTicker(ui *UI) {
	if ui.isTickerRunning.Load() && ui.tickerDone != nil {
		ui.tickerDone <- true
	}
}
