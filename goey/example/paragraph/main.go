// This package provides an example application built using the goey package
// that demonstrates using the Image widget.  Clicking the button will cycle
// through images that are a uniform colour, as well as an image of the Go
// mascot.
//
// The management of scrollbars can be tested by using the environment variable
// GOEY_SCROLL.  Allowed values are 0 through 3, which enable no scrollbars,
// the vertical scrollbar, the horizontal scrollbar, or both scrollbars.
package main

import (
	"fmt"
	_ "image/png"

	"guitest/goey"
	"guitest/goey/base"
	"guitest/goey/loop"
)

var (
	mainWindow *goey.Window
	align      goey.TextAlignment
)

const (
	lorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum"
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	// Add the controls
	mw, err := goey.NewWindow("Paragraph", render())
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
}

func render() base.Widget {
	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			AlignMain:  goey.MainCenter,
			AlignCross: goey.CrossCenter,
			Children: []base.Widget{
				&goey.P{Text: lorem, Align: align},
				&goey.Expand{},
				&goey.HR{},
				&goey.SelectInput{
					Items:    []string{"Left align", "Center align", "Right align", "Justify"},
					Value:    int(align),
					OnChange: onChangeAlign,
				},
			},
		},
	}
}

func onChangeAlign(value int) {
	align = goey.TextAlignment(value)
	update()
}
