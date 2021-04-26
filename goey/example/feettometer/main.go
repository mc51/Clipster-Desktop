// This package provides an example application built using the goey package
// that rebuilds the classic Tcl/Tk tutorial application.
//
// The example also shows the use of a custom layout container, MinSizedBox,
// showing that new layouts can be developed outside of the main package,
// and used portably.  In this case, the layout uses the methods
// MinIntrinsicHeight and MinIntrinsicWidth to find the minimum acceptable size
// for the child, and then limits the child to that particular size as long as
// it meets the layout constraints.
//
// The GUI is partially dynamic, in that the conversion from feet to meters is
// performed whenever the button is pressed.  However, it would be very easy
// to have the conversion performed continuously as the user types by adding
// a call to the conversion function included in the OnChange callback of the
// textbox.
//
// The management of scrollbars can be tested by using the environment variable
// GOEY_SCROLL.  Allowed values are 0 through 3, which enable no scrollbars,
// the vertical scrollbar, the horizontal scrollbar, or both scrollbars.
package main

import (
	"fmt"
	"strconv"

	"clipster/goey"
	"clipster/goey/base"
	"clipster/goey/loop"
)

var (
	mainWindow *goey.Window

	feetValue  string
	meterValue string
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func createWindow() error {
	// Add the controls
	mw, err := goey.NewWindow("Feet to Meters", render())
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
		Child: &goey.Align{Child: &MinSizedBox{Child: &goey.VBox{
			AlignMain: goey.MainCenter,
			Children: []base.Widget{
				&goey.HBox{
					AlignMain:  goey.Homogeneous,
					AlignCross: goey.CrossCenter,
					Children: []base.Widget{
						&goey.Empty{},
						&goey.TextInput{Value: feetValue, OnChange: func(v string) { feetValue = v }, OnEnterKey: func(v string) { feetValue = v; calculate() }},
						&goey.Label{Text: "feet"},
					},
				}, &goey.HBox{
					AlignMain:  goey.Homogeneous,
					AlignCross: goey.CrossCenter,
					Children: []base.Widget{
						&goey.Label{Text: "is equivalent to"},
						&goey.Label{Text: meterValue},
						&goey.Label{Text: "meters"},
					},
				}, &goey.HBox{
					AlignMain:  goey.Homogeneous,
					AlignCross: goey.CrossCenter,
					Children: []base.Widget{
						&goey.Empty{},
						&goey.Empty{},
						&goey.Button{Text: "Calculate", Default: true, OnClick: calculate},
					},
				},
			},
		}}},
	}
}

func calculate() {
	feet, err := strconv.ParseFloat(feetValue, 64)
	if err != nil {
		meterValue = "(error)"
	} else {
		meterValue = fmt.Sprintf("%f", feet*0.3048)
	}
	update()
}
