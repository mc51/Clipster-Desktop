package goey

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"testing/quick"

	"guitest/goey/base"
)

func ExampleButton() {
	clickCount := 0

	// In a full application, this variable would be updated to point to
	// the main window for the application.
	var mainWindow *Window
	// These functions are used to update the GUI.  See below
	var update func()
	var render func() base.Widget

	// Update function
	update = func() {
		err := mainWindow.SetChild(render())
		if err != nil {
			panic(err)
		}
	}

	// Render function generates a tree of Widgets to describe the desired
	// state of the GUI.
	render = func() base.Widget {
		// Prep - text for the button
		text := "Click me!"
		if clickCount > 0 {
			text = text + "  (" + strconv.Itoa(clickCount) + ")"
		}
		// The GUI contains a single widget, this button.
		return &VBox{
			AlignMain:  MainCenter,
			AlignCross: CrossCenter,
			Children: []base.Widget{
				&Button{Text: text, OnClick: func() {
					clickCount++
					update()
				}},
			}}
	}
}

func buttonValues(values []reflect.Value, rand *rand.Rand) {
	// Get a string
	labelValues(values, rand)

	// Create a choices for disabled and default
	values[1] = reflect.ValueOf(rand.Uint64()%2 == 0)
	values[2] = reflect.ValueOf(rand.Uint64()%2 == 0)
}

func TestButtonMount(t *testing.T) {
	testingMountWidgets(t,
		&Button{Text: "A"},
		&Button{Text: "D", Disabled: true},
		&Button{Text: "E", Default: true},
	)

	t.Run("QuickCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		f := func(text string, disabled, def bool) bool {
			return testingMountWidget(t, &Button{Text: text, Disabled: disabled, Default: def})
		}
		if err := quick.Check(f, &quick.Config{Values: buttonValues}); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}

func TestButtonClose(t *testing.T) {
	testingCloseWidgets(t,
		&Button{Text: "A"},
		&Button{Text: "D", Disabled: true},
		&Button{Text: "E", Default: true},
	)
}

func TestButtonFocus(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	)
}

func TestButtonClick(t *testing.T) {
	testingCheckClick(t,
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	)
}

func TestButtonUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Button{Text: "A"},
		&Button{Text: "D", Disabled: true},
		&Button{Text: "E", Default: true},
	}, []base.Widget{
		&Button{Text: "AB"},
		&Button{Text: "DB", Default: true},
		&Button{Text: "EB", Disabled: true},
	})

	t.Run("QuickCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		updater, closer := testingUpdateWidget(t)
		defer closer()

		f := func(text string, disabled, def bool) bool {
			return updater(&Button{Text: text, Disabled: disabled, Default: def})
		}
		if err := quick.Check(f, &quick.Config{Values: buttonValues}); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}
