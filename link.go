package main

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"os/exec"
	"runtime"
)

type Hyperlink struct {
	widget.Clickable
	Text string
	URL  string
}

// openURL opens a URL in the default browser
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// Layout implements the custom hyperlink layout
func (h *Hyperlink) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Color for the link
	linkColor := color.NRGBA{R: 0, G: 0, B: 238, A: 255} // Standard blue link color

	// Create the label
	label := material.Body1(th, h.Text)
	label.Color = linkColor

	// Handle click

	if h.Clickable.Clicked(gtx) {
		go func() {
			err := openURL(h.URL)
			if err != nil {
				println("Error opening URL:", err.Error())
			}
		}()
	}

	return h.Clickable.Layout(gtx, label.Layout)
}
