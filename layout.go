package main

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
)

func drawImage(gtx layout.Context, img paint.ImageOp, originalSize image.Point) layout.Dimensions {
	// Get available space from constraints
	availWidth := float32(gtx.Constraints.Max.X)
	availHeight := float32(gtx.Constraints.Max.Y)

	// Get original image dimensions
	imgWidth := float32(originalSize.X)
	imgHeight := float32(originalSize.Y)

	// Calculate scale to fit in window while maintaining aspect ratio
	scale := minC(availWidth/imgWidth, availHeight/imgHeight)

	// Calculate new dimensions
	newWidth := int(imgWidth * scale)
	newHeight := int(imgHeight * scale)

	// Center the image
	offsetX := (int(availWidth) - newWidth) / 2
	offsetY := (int(availHeight) - newHeight) / 2

	// Create a stack for transformations
	macro := op.Record(gtx.Ops)

	// First, apply the offset for centering
	op.Offset(image.Pt(offsetX, offsetY)).Add(gtx.Ops)

	// Then apply scaling
	scaleMacro := op.Record(gtx.Ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Point{}, f32.Pt(scale, scale))).Add(gtx.Ops)

	// Create clip rect for the original size (before scaling)
	clip.Rect{Max: originalSize}.Push(gtx.Ops).Pop()

	// Draw the image
	img.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	scaleMacro.Stop().Add(gtx.Ops)
	macro.Stop().Add(gtx.Ops)

	// Return dimensions using the full available space
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func getImg(ui *UI) IMG {
	if ui.mainImage == 1 {
		return ui.img1
	}
	return ui.img2
}

// TODO make this more elegant
func createLayout(gtx layout.Context,

	th *material.Theme, startButton, stopButton, setButton, aboutButton *widget.Clickable, ui *UI) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Spacing:   layout.SpaceBetween,
		Alignment: layout.Middle, // Changed to Middle for better centering
	}.Layout(gtx,
		// Image container that takes all available space
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if ui.isTickerRunning.Load() {
				// Pass the full context constraints to drawImage
				return drawImage(gtx, getImg(ui).imgOp, getImg(ui).imgSize)
			}
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		// Button container
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Spacing:   layout.SpaceEvenly,
					Alignment: layout.Middle,
				}.Layout(gtx,
					// Label
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(th, "Flicker rate (per second)")
						label.Alignment = text.Middle
						return label.Layout(gtx)
					}),
					//space between
					layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
					// Rate Editor
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(60)
						editor := material.Editor(th, &ui.rateEditor, "Rate")
						return editor.Layout(gtx)
					}),

					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(100)
						gtx.Constraints.Min.Y = gtx.Dp(50)
						return createButton(gtx, th, startButton, "Start")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(100)
						gtx.Constraints.Min.Y = gtx.Dp(50)
						return createButton(gtx, th, stopButton, "Stop")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(100)
						gtx.Constraints.Min.Y = gtx.Dp(50)
						return createButton(gtx, th, setButton, "Set")
					}),

					layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),

					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(100)
						gtx.Constraints.Min.Y = gtx.Dp(50)
						return createButton(gtx, th, aboutButton, "About")
					}),
				)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
	)
}
func createButton(gtx layout.Context, th *material.Theme, buttonWidget *widget.Clickable, buttonText string) layout.Dimensions {
	btn := material.Button(th, buttonWidget, buttonText)
	return btn.Layout(gtx)
}

func minC(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
