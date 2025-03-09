package main

import (
	"fmt"
	"gioui.org/app"
	"os"
	"strings"
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

	// Check if we're using schedule mode
	if ui.useSchedule && len(ui.schedule) > 0 {
		startScheduledTicker(ui, w)
		return
	}

	// Normal mode (single rate)
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

func startScheduledTicker(ui *UI, w *app.Window) {
	// Initialize schedule variables
	ui.currentScheduleIndex = 0
	ui.scheduleStartTime = time.Now()

	// Start with the first schedule item
	currentItem := ui.schedule[0]

	var rate int
	if currentItem.BlankTime > 0 {
		// This is a blank period, we don't need to flicker
		rate = 1000 // Just a slow tick to check for schedule transitions
	} else {
		var err error
		rate, err = calculateRate(currentItem.FlickeringRate)
		if err != nil {
			rate = 1000
		}
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

		// Create a ticker to check schedule transitions every 100ms
		scheduleTicker := time.NewTicker(100 * time.Millisecond)
		defer scheduleTicker.Stop()

		for {
			select {
			case <-ui.t.C:
				// Only flip if we're not in a blank period
				currentItem := ui.schedule[ui.currentScheduleIndex]
				if currentItem.BlankTime == 0 {
					ui.FlipImage(w)
				}

			case <-scheduleTicker.C:
				// Check if we need to transition to the next schedule item
				elapsed := time.Since(ui.scheduleStartTime).Seconds()
				currentItemDuration := 0

				// Calculate the total duration up to the current item
				totalDuration := 0
				for i := 0; i <= ui.currentScheduleIndex; i++ {
					currentItemDuration = ui.schedule[i].Duration
					totalDuration += currentItemDuration
				}

				// Check if we need to move to the next item
				if elapsed >= float64(totalDuration) {
					ui.currentScheduleIndex++

					// If we've reached the end of the schedule, start over
					if ui.currentScheduleIndex >= len(ui.schedule) {
						ui.currentScheduleIndex = 0
						ui.scheduleStartTime = time.Now()
					}

					// Update the ticker rate for the new schedule item
					newItem := ui.schedule[ui.currentScheduleIndex]
					var newRate int

					if newItem.BlankTime > 0 {
						// This is a blank period
						newRate = 1000 // Slow tick rate for blank periods
					} else {
						var err error
						newRate, err = calculateRate(newItem.FlickeringRate)
						if err != nil {
							newRate = 1000
						}
					}

					// Update the ticker
					ui.t.Stop()
					ui.t = time.NewTicker(time.Millisecond * time.Duration(newRate))
				}

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

// Parse schedule text into ScheduleItem structs
func parseSchedule(ui *UI, scheduleText string) {
	ui.schedule = []ScheduleItem{}

	// Split by semicolons
	parts := []string{}
	for _, part := range strings.Split(scheduleText, ";") {
		part = strings.TrimSpace(part)
		if part != "" {
			parts = append(parts, part)
		}
	}

	for _, part := range parts {
		// Check if it's a blank screen time or a flickering period
		if !strings.Contains(part, "-") {
			// This is a blank screen time
			var blankTime int
			_, err := fmt.Sscanf(part, "%d", &blankTime)
			if err == nil && blankTime > 0 {
				ui.schedule = append(ui.schedule, ScheduleItem{
					Duration:       blankTime,
					FlickeringRate: 0,
					BlankTime:      blankTime,
				})
			}
		} else {
			// This is a flickering period
			var duration, rate int
			_, err := fmt.Sscanf(part, "%d-%d", &duration, &rate)
			if err == nil && duration > 0 && rate > 0 {
				ui.schedule = append(ui.schedule, ScheduleItem{
					Duration:       duration,
					FlickeringRate: rate,
					BlankTime:      0,
				})
			}
		}
	}
}

// Save schedule to file
func saveSchedule(ui *UI) {
	scheduleText := ui.scheduleEditor.Text()
	err := os.WriteFile("schedule.txt", []byte(scheduleText), 0644)
	if err != nil {
		// Handle error (could show in UI but for now we'll just ignore)
		fmt.Println("Error saving schedule:", err)
	}
}
