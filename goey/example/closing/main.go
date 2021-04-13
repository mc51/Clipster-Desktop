// This package provides an example application built using the goey package
// that demonstrates using the OnClosing callback for windows.  Trying to close
// the window using the normal method will fail, but the button within the
// window can be used.
package main

import (
	"fmt"

	"guitest/goey"
	"guitest/goey/base"
	"guitest/goey/loop"
)

var (
	mainWindow *goey.Window
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	mw, err := goey.NewWindow("Closing", render())
	if err != nil {
		return err
	}
	mw.SetOnClosing(func() bool {
		// Block closing of the window
		return true
	})
	mainWindow = mw

	return nil
}

func render() base.Widget {
	return &goey.Padding{
		Insets: goey.UniformInsets(36 * goey.DIP),
		Child: &goey.Align{
			Child: &goey.Button{Text: "Close app", OnClick: func() {
				mainWindow.Close()
			}},
		},
	}
}
