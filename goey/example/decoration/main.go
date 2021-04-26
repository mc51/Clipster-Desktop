// This package provides an example application built using the goey package
// that demontrates using the Decoration widget.  Clicking the button will cycle
// through various background colours, occasionally adding a border, and also
// cycle through increasing the border radius.
//
// The management of scrollbars can be tested by using the environment variable
// GOEY_SCROLL.  Allowed values are 0 through 3, which enable no scrollbars,
// the vertical scrollbar, the horizontal scrollbar, or both scrollbars.
package main

import (
	"fmt"
	"image/color"
	"strconv"

	"clipster/goey"
	"clipster/goey/base"
	"clipster/goey/loop"
)

var (
	mainWindow *goey.Window
	clickCount int

	colors = [4]color.RGBA{
		{0xc0, 0x40, 0x40, 0xff},
		{0x40, 0xc0, 0x40, 0xff},
		{0x40, 0x40, 0xc0, 0xff},
		{0, 0, 0, 0},
	}
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	mw, err := goey.NewWindow("Decoration", render())
	if err != nil {
		return err
	}
	mainWindow = mw

	return nil
}

func updateWindow() {
	err := mainWindow.SetChild(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func render() base.Widget {
	text := "Click me!"
	if clickCount > 0 {
		text = text + "  (" + strconv.Itoa(clickCount) + ")"
	}
	stroke := color.RGBA{}
	if clickCount > 0 && clickCount%3 == 0 {
		stroke.A = 0xFF
	}
	return &goey.VBox{
		AlignMain:  goey.SpaceAround,
		AlignCross: goey.CrossCenter,
		Children: []base.Widget{
			&goey.Padding{
				Insets: goey.DefaultInsets(),
				Child: &goey.P{
					Text: "This is a demonstration of the use of a Decoration widget.  Clicking the button will cycle through different background colours, as well as change the border radius.",
				},
			},
			&goey.Decoration{
				Fill:   colors[clickCount%4],
				Stroke: stroke,
				Insets: goey.UniformInsets(0.5 * 96 * goey.DIP),
				Radius: (2 * goey.DIP).Scale(clickCount%16, 1),
				Child: &goey.Button{Text: text, OnClick: func() {
					clickCount++
					updateWindow()
				}},
			},
		},
	}
}
