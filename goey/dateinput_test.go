package goey

import (
	"testing"
	"time"

	"clipster/goey/base"
)

func TestDateInputMount(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 0, 0, 0, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 0, 0, 0, 0, time.Local)

	testingMountWidgets(t,
		&DateInput{Value: v1},
		&DateInput{Value: v2, Disabled: true},
		&DateInput{Value: v2},
	)
}

func TestDateInputClose(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 15, 4, 5, 0, time.Local)

	testingCloseWidgets(t,
		&DateInput{Value: v1},
		&DateInput{Value: v2, Disabled: true},
		&DateInput{Value: v2},
	)
}

func TestDateInputEvents(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 15, 4, 5, 0, time.Local)

	testingCheckFocusAndBlur(t,
		&DateInput{Value: v1},
		&DateInput{Value: v2},
		&DateInput{Value: v2},
	)
}

func TestDateInputUpdateProps(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 0, 0, 0, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 0, 0, 0, 0, time.Local)

	testingUpdateWidgets(t, []base.Widget{
		&DateInput{Value: v1},
		&DateInput{Value: v2, Disabled: true},
		&DateInput{Value: v2},
	}, []base.Widget{
		&DateInput{Value: v2},
		&DateInput{Value: v2, Disabled: false},
		&DateInput{Value: v1, Disabled: true},
	})
}
