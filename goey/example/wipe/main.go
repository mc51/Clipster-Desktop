// This package provides an example application built using the goey package
// that demonstrates using animation.  Clicking the buttons will cycle
// through images that are a uniform colour, performing a slide as the image
// is updated.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"guitest/goey"
	"guitest/goey/animate"
	"guitest/goey/base"
	"guitest/goey/loop"
)

var (
	mainWindow *goey.Window
	clickCount int

	colors = [3]color.RGBA{
		{0xff, 0, 0, 0xff},
		{0, 0xff, 0, 0xff},
		{0, 0, 0xff, 0xff},
	}
	colorNames = [3]string{
		"Red", "Green", "Blue",
	}
)

func selectImage(index int) (image.Image, string) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(img, img.Bounds(), image.NewUniform(colors[index%3]), image.Point{}, draw.Src)
	return img, colorNames[index%3]
}

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	// Add the controls
	mw, err := goey.NewWindow("Wipe", render())
	if err != nil {
		return err
	}
	mainWindow = mw

	return nil
}

func update() {
	err := mainWindow.SetChild(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	img, _ := selectImage(clickCount)
	mainWindow.SetIcon(img)
}

func render() base.Widget {
	img, description := selectImage(clickCount)

	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			AlignMain:  goey.MainCenter,
			AlignCross: goey.Stretch,
			Children: []base.Widget{
				&goey.HBox{
					AlignMain: goey.MainCenter,
					Children: []base.Widget{
						&goey.Button{Text: "Increment", OnClick: func() {
							clickCount++
							update()
						}},
						&goey.Button{Text: "Decrement", Disabled: clickCount == 0, OnClick: func() {
							clickCount--
							update()
						}},
					},
				},
				&animate.Wipe{
					Child: &goey.Align{
						Child: &goey.Img{
							Image:  img,
							Width:  (1 * goey.DIP).Scale(img.Bounds().Dx(), 1),
							Height: (1 * goey.DIP).Scale(img.Bounds().Dy(), 1),
						},
					},
					Level: clickCount,
				},
				&goey.P{Text: description, Align: goey.JustifyCenter},
			},
		},
	}
}
