// This package provides an example application built using the goey package
// that shows a sidebar an array of buttons.  This is meant to be an example
// of a corporate portal.
package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"

	"guitest/goey"
	"guitest/goey/base"
	"guitest/goey/loop"
)

var (
	gopher image.Image
	window *goey.Window
)

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close of a file with only read permission.  Will not error.
		_ = file.Close()
	}()

	img, _, err := image.Decode(file)
	return img, err
}

func main() {
	img, err := loadImage("gopher.png")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	gopher = img

	err = loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	w, err := goey.NewWindow("Menu", renderWindow())
	if err != nil {
		return err
	}
	window = w
	return nil
}

func updateWindow() {
	window.SetChild(renderWindow())
}

func renderSidebar() base.Widget {
	return &goey.Decoration{
		Fill: color.RGBA{128, 255, 128, 255},
		Child: &goey.Padding{
			Insets: goey.DefaultInsets(),
			Child: &goey.VBox{goey.MainCenter, goey.CrossCenter, []base.Widget{
				&goey.Label{Text: "Example Menu"},
				&goey.Img{Image: gopher},
			},
			},
		}}
}

func renderMainbar() base.Widget {
	return &goey.Expand{Child: &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{goey.MainCenter, goey.Stretch, []base.Widget{
			&Column{[]base.Widget{
				&goey.Button{Text: "A1"}, &goey.Button{Text: "A2"}, &goey.Button{Text: "A3"}, &goey.Button{Text: "A4"},
				&goey.Button{Text: "B1"}, &goey.Button{Text: "B2"}, &goey.Button{Text: "B3"}, &goey.Button{Text: "B4"},
				&goey.Button{Text: "C1"}, &goey.Button{Text: "C2"}, &goey.Button{Text: "C3"}, &goey.Button{Text: "C4"},
				&goey.Button{Text: "D1"}, &goey.Button{Text: "D2"}, &goey.Button{Text: "D3"}, &goey.Button{Text: "D4"},
			}},
			&goey.HR{},
			&goey.Button{Text: "Help"},
		}},
	}}
}

func renderWindow() base.Widget {
	ret := &goey.HBox{
		Children: []base.Widget{
			renderSidebar(),
			renderMainbar(),
		},
	}

	return ret
}
