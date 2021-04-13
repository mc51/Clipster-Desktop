// This package provides an example application built using the goey package
// that shows use of the Align layout widget.
package main

import (
	"fmt"

	"guitest/goey"
	"guitest/goey/base"
	"guitest/goey/loop"
)

var (
	mainWindow *goey.Window
	halign     goey.Alignment
	valign     goey.Alignment
	wf         float64 = 2
	hf         float64 = 2
)

func main() {
	init := func() error {
		mw, err := goey.NewWindow("Align Widget Example", render())
		if err != nil {
			return err
		}

		mainWindow = mw
		return nil
	}

	err := loop.Run(init)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func updateWindow() {
	// To update the window, we generate a new widget for the contents of the
	// top-level window.
	err := mainWindow.SetChild(render())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func render() base.Widget {
	return &goey.Padding{
		Child: &goey.VBox{
			Children: []base.Widget{
				&goey.HBox{
					AlignMain: goey.Homogeneous,
					Children: []base.Widget{
						AlignmentInput(&halign),
						AlignmentInput(&valign),
					},
				},
				&goey.HR{},
				&goey.Expand{
					Child: &goey.Align{
						HAlign: halign,
						VAlign: valign,
						Child:  &goey.Button{Text: "Noop Button"},
					},
				},
			},
		},
		Insets: goey.DefaultInsets(),
	}
}

func AlignmentInput(value *goey.Alignment) base.Widget {
	ndx := func(a goey.Alignment) int {
		switch a {
		case goey.AlignStart:
			return 0
		case goey.AlignCenter:
			return 1
		default:
			return 2
		}
	}(*value)

	return &goey.SelectInput{
		Items: []string{"Start", "Center", "End"},
		Value: ndx,
		OnChange: func(newValue int) {
			*value = func(a int) goey.Alignment {
				switch a {
				case 0:
					return goey.AlignStart
				case 1:
					return goey.AlignCenter
				default:
					return goey.AlignEnd
				}
			}(newValue)
			updateWindow()
		},
	}
}
