package main

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
)

type AboutDialog struct {
	isOpen      bool
	closeButton widget.Clickable
	websiteLink Hyperlink
	githubLink  Hyperlink
}

func NewAboutDialog() *AboutDialog {
	return &AboutDialog{
		websiteLink: Hyperlink{
			Text: "Some notes about this app and studies that inspired it",
			URL:  "https://2bits2bits2.github.io/2bits2bits2notes/go-brain-flicker",
		},
		githubLink: Hyperlink{
			Text: "GitHub repository",
			URL:  "https://github.com/2bits2bits2/go-brain-flicker",
		},
	}
}

func (d *AboutDialog) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if !d.isOpen {
		return layout.Dimensions{}
	}

	gtx.Constraints.Min = image.Point{X: gtx.Dp(300), Y: gtx.Dp(400)}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(d.layoutBackground),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(d.layoutPopupBackground(gtx)),
					layout.Stacked(d.layoutContent(th)),
				)
			})
		}),
	)
}

func (d *AboutDialog) layoutBackground(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, color.NRGBA{A: 200})
	return layout.Dimensions{Size: gtx.Constraints.Min}
}

func (d *AboutDialog) layoutPopupBackground(gtx layout.Context) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		rr := clip.RRect{
			Rect: image.Rectangle{Max: gtx.Constraints.Min},
			SE:   gtx.Dp(8), SW: gtx.Dp(8),
			NE: gtx.Dp(8), NW: gtx.Dp(8),
		}
		paint.FillShape(gtx.Ops,
			color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			clip.RRect{
				Rect: rr.Rect,
				NE:   rr.NE, NW: rr.NW,
				SE: rr.SE, SW: rr.SW,
			}.Op(gtx.Ops))
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
}

func (d *AboutDialog) layoutContent(th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					d.layoutTitle(th),
					space(16),
					d.layoutBody(th),
					space(16),
					d.layoutCloseButton(th),
				)
			},
		)
	}
}

func (d *AboutDialog) layoutTitle(th *material.Theme) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		title := material.H6(th, "About Brain Flicker")
		title.Alignment = text.Middle
		return title.Layout(gtx)
	})
}

func (d *AboutDialog) layoutBody(th *material.Theme) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th, "Brain Flicker is an application for visual experiments.").Layout(gtx)
			}),
			space(8),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.websiteLink.Layout(gtx, th)
			}),
			space(8),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.githubLink.Layout(gtx, th)
			}),
			space(8),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th, "Version: 0.1").Layout(gtx)
			}),
		)
	})
}

func (d *AboutDialog) layoutCloseButton(th *material.Theme) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		btn := material.Button(th, &d.closeButton, "Close")
		return btn.Layout(gtx)
	})
}

func space(size int) layout.FlexChild {
	return layout.Rigid(layout.Spacer{Height: unit.Dp(size)}.Layout)
}
