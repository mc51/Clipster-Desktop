package goey

import (
	"clipster/goey/base"
)

var (
	intInputKind = base.NewKind("clipster/goey.IntInput")
)

// IntInput describes a widget that users input or update a single integer value.
// The model for the value is a int64.
//
// If the field Min and Max are both zero, then a default range will be
// initialized covering the entire range of int64.  Note that on some platforms,
// the control internally is a float64, so the entire range of int64 cannot be
// covered without lose of precision.
type IntInput struct {
	Value       int64             // Value is the current value for the field
	Placeholder string            // Placeholder is a descriptive text that can be displayed when the field is empty
	Disabled    bool              // Disabled is a flag indicating that the user cannot interact with this field
	Min, Max    int64             // Min and Max set the range of Value
	OnChange    func(value int64) // OnChange will be called whenever the user changes the value for this field
	OnFocus     func()            // OnFocus will be called whenever the field receives the keyboard focus
	OnBlur      func()            // OnBlur will be called whenever the field loses the keyboard focus
	OnEnterKey  func(value int64) // OnEnterKey will be called whenever the use hits the enter key
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*IntInput) Kind() *base.Kind {
	return &intInputKind
}

// Mount creates a text field in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *IntInput) Mount(parent base.Control) (base.Element, error) {
	// Fill in default values for the range.
	w.UpdateRange()
	// Make sure that the value is within the range.
	w.UpdateValue()

	// Forward to the platform-dependant code
	return w.mount(parent)
}

// UpdateRange sets a default range when the fields Min and Max are both
// default initialized.  The default range matches the range of int64.
func (w *IntInput) UpdateRange() {
	if w.Min == 0 && w.Max == 0 {
		// See document for package builtin, type int64
		w.Min = -9223372036854775808
		w.Max = 9223372036854775807
	}
}

// UpdateValue clamps the field Value to the range [Min,Max].
func (w *IntInput) UpdateValue() {
	if w.Value < w.Min {
		w.Value = w.Min
	} else if w.Value > w.Max {
		w.Value = w.Max
	}
}

func (*intinputElement) Kind() *base.Kind {
	return &intInputKind
}

func (w *intinputElement) UpdateProps(data base.Widget) error {
	widget := data.(*IntInput)

	// Fill in default values for the range.
	widget.UpdateRange()
	widget.UpdateValue()
	// Forward to the platform-dependant code
	return w.updateProps(widget)
}
