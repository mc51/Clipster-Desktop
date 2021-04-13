package goey

import (
	"guitest/goey/base"
)

var (
	progressKind = base.NewKind("guitest/goey.Progress")
)

// Progress describes a widget that shows a progress bar.
// The model for the value is an int.
//
// If both Min and Max are zero, then Max will be updated to 100.  Other cases
// where Min == Max are not allowed.
type Progress struct {
	Value    int // Value is the current value to be displayed
	Min, Max int // Min and Max set the range of Value
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Progress) Kind() *base.Kind {
	return &progressKind
}

// Mount creates a progress control in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Progress) Mount(parent base.Control) (base.Element, error) {
	// Fill in default values for the range.
	w.UpdateRange()
	// Make sure that the value is within the range.
	w.UpdateValue()

	// Forward to the platform-dependant code
	return w.mount(parent)
}

// UpdateRange sets a default range when the fields Min and Max are both
// default initialized.
func (w *Progress) UpdateRange() {
	// Fill in a default range if none has been specified.
	if w.Min == 0 && w.Max == 0 {
		w.Max = 100
	}
}

// UpdateValue clamps the field Value to the range [Min,Max].
func (w *Progress) UpdateValue() {
	if w.Value < w.Min {
		w.Value = w.Min
	} else if w.Value > w.Max {
		w.Value = w.Max
	}
}

func (*progressElement) Kind() *base.Kind {
	return &progressKind
}

func (w *progressElement) UpdateProps(data base.Widget) error {
	pb := data.(*Progress)

	// Fill in default values for the range.
	pb.UpdateRange()
	pb.UpdateValue()
	// Forward to the platform-dependant code
	return w.updateProps(pb)
}
