// This package provides an example application built using the goey package
// that shows the use of message boxes.
package main

import (
	"fmt"

	"guitest/goey"
	"guitest/goey/base"
	"guitest/goey/dialog"
	"guitest/goey/loop"
)

var (
	window *goey.Window
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	// This is the callback used to initialize the GUI state.  For this simple
	// example, we need to create a new top-level window, and set a child
	// widget.
	mw, err := goey.NewWindow("Message Box", render())
	if err != nil {
		return err
	}

	// Stash a copy of the top-level window for use with call to message box.
	window = mw

	return nil
}

func render() base.Widget {
	// We return a widget describing the desired state of the GUI.  Note that
	// this is data only, and no changes have been effected yet.
	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			AlignMain: goey.MainCenter,
			Children: []base.Widget{
				&goey.Button{
					Text: "Show simple message box",
					OnClick: func() {
						window.Message("This is the body of the message").Show()
					},
				},
				&goey.Button{
					Text: "Show complete message box",
					OnClick: func() {
						window.Message("This is the body of the message").WithInfo().WithTitle("Custom title").Show()
					},
				},
				&goey.Button{
					Text: "Show message box without parent",
					OnClick: func() {
						dialog.NewMessage("This is the body of the message.  Note, you can still interact with the main window.").Show()
					},
				},
			},
		},
	}
}
