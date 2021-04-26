package goey

import (
	"clipster/goey/base"
	"testing"
)

func TestSelectInputMount(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testingMountWidgets(t,
		&SelectInput{Value: 0, Items: options},
		&SelectInput{Value: 1, Items: options},
		&SelectInput{Value: 2, Items: options, Disabled: true},
		&SelectInput{Unset: true, Items: options, Disabled: true},
		&SelectInput{Items: []string{}},
	)
}

func TestSelectInputClose(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testingCloseWidgets(t,
		&SelectInput{Value: 0, Items: options},
		&SelectInput{Value: 1, Items: options},
		&SelectInput{Value: 2, Items: options, Disabled: true},
		&SelectInput{Unset: true, Items: options, Disabled: true},
	)
}

func TestSelectInputEvents(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testingCheckFocusAndBlur(t,
		&SelectInput{Items: options},
		&SelectInput{Items: options},
		&SelectInput{Items: options},
	)
}

func TestSelectInputUpdateProps(t *testing.T) {
	options1 := []string{"Option A", "Option B", "Option C"}
	options2 := []string{"Choice A", "Choice B", "Choice C"}

	testingUpdateWidgets(t, []base.Widget{
		&SelectInput{Value: 0, Items: options1},
		&SelectInput{Value: 1, Items: options2},
		&SelectInput{Value: 2, Items: options1, Disabled: true},
		&SelectInput{Unset: true, Items: options2},
	}, []base.Widget{
		&SelectInput{Value: 1, Items: options2},
		&SelectInput{Unset: true, Items: options1},
		&SelectInput{Value: 2, Items: options1, Disabled: true},
		&SelectInput{Value: 1, Items: options2},
	})
}
