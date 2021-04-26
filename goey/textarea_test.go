package goey

import (
	"bytes"
	"testing"

	"clipster/goey/base"
)

func TestTextAreaMount(t *testing.T) {
	// Note, cannot use zero value for MinLines.  This will be changed to a
	// default value, and cause the post mounting check that the widget was
	// correctly instantiated to fail.
	testingMountWidgets(t,
		&TextArea{Value: "A", MinLines: 3},
		&TextArea{Value: "B", MinLines: 3, Placeholder: "..."},
		&TextArea{Value: "C", MinLines: 3, Disabled: true},
	)
}

func TestTextAreaOnFocus(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&TextArea{},
		&TextArea{},
		&TextArea{},
	)
}

func TestTextAreaOnChange(t *testing.T) {
	log := bytes.NewBuffer(nil)

	testingTypeKeys(t, "Hello",
		&TextArea{OnChange: func(v string) {
			log.WriteString(v)
			log.WriteString("\x1E")
		}})

	const want = "H\x1EHe\x1EHel\x1EHell\x1EHello\x1E"
	if got := log.String(); got != want {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func TestTextAreaUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&TextArea{Value: "A", MinLines: 5},
		&TextArea{Value: "B", MinLines: 3, Placeholder: "..."},
		&TextArea{Value: "C", MinLines: 3, Disabled: true},
	}, []base.Widget{
		&TextArea{Value: "AA", MinLines: 6},
		&TextArea{Value: "BA", MinLines: 3, Disabled: true},
		&TextArea{Value: "CA", MinLines: 3, Placeholder: "***", Disabled: false},
	})
}
