package goey

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"clipster/goey/base"
)

func checkboxValues(values []reflect.Value, rand *rand.Rand) {
	// Get a string
	labelValues(values, rand)

	// Create a choices for value and disabled
	values[1] = reflect.ValueOf(rand.Uint64()%2 == 0)
	values[2] = reflect.ValueOf(rand.Uint64()%2 == 0)
}

func TestCheckboxMount(t *testing.T) {
	testingMountWidgets(t,
		&Checkbox{Value: false, Text: "A"},
		&Checkbox{Value: true, Text: "B"},
		&Checkbox{Value: false, Text: "C", Disabled: true},
		&Checkbox{Value: true, Text: "D", Disabled: true},
		&Checkbox{Text: ""},
	)

	t.Run("QuickCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		f := func(text string, value, disabled bool) bool {
			return testingMountWidget(t, &Checkbox{Text: text, Value: value, Disabled: disabled})
		}
		if err := quick.Check(f, &quick.Config{Values: checkboxValues}); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}

func TestCheckboxClose(t *testing.T) {
	testingCloseWidgets(t,
		&Checkbox{Value: false, Text: "A"},
		&Checkbox{Value: true, Text: "B", Disabled: true},
	)
}

func TestCheckboxFocus(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&Checkbox{Text: "A"},
		&Checkbox{Text: "B"},
		&Checkbox{Text: "C"},
	)
}

func TestCheckboxClick(t *testing.T) {
	var values [3]bool

	testingCheckClick(t,
		&Checkbox{Text: "A", OnChange: func(v bool) { values[0] = v }},
		&Checkbox{Text: "B", Value: true, OnChange: func(v bool) { values[1] = v }},
		&Checkbox{Text: "C", OnChange: func(v bool) { values[2] = v }},
	)

	if !values[0] || values[1] || !values[2] {
		t.Errorf("OnChange failed, expected %v, got %v", [3]bool{true, false, true}, values[:])
	}
}

func TestCheckboxUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Checkbox{Value: false, Text: "A"},
		&Checkbox{Value: true, Text: "B", Disabled: true},
	}, []base.Widget{
		&Checkbox{Value: true, Text: "A--", Disabled: true},
		&Checkbox{Value: false, Text: "B--", Disabled: false},
	})
}
