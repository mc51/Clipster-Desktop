package goey

import (
	"guitest/goey/base"
)

var (
	sliderKind = base.NewKind("guitest/goey.Slider")
)

// Slider describes a widget that users input or update a single real value.
// The model for the value is a float64.
//
// If both Min and Max are zero, then Max will be updated to 100.  Other cases
// where Min == Max are not allowed.
type Slider struct {
	Value    float64       // Value is the current value for the field
	Disabled bool          // Disabled is a flag indicating that the user cannot interact with this field
	Min, Max float64       // Min and Max set the range of Value
	OnChange func(float64) // OnChange will be called whenever the user changes the value for this field
	OnFocus  func()        // OnFocus will be called whenever the slider receives the keyboard focus.
	OnBlur   func()        // OnBlur will be called whenever the slider loses the keyboard focus.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Slider) Kind() *base.Kind {
	return &sliderKind
}

// Mount creates a slider control in the GUI.
// The newly created widget will be a child of the widget specified by parent.
//
// If both fields Min and Max are zero, a default range will be created going
// from 0 to 100.
//
// The field Value will be updated to lie within the range.
func (w *Slider) Mount(parent base.Control) (base.Element, error) {
	// Fill in default values for the range.
	w.UpdateRange()
	// Make sure that the value is within the range.
	w.UpdateValue()

	// Forward to the platform-dependant code
	return w.mount(parent)
}

// UpdateRange sets a default range when the fields Min and Max are both
// default initialized.
func (w *Slider) UpdateRange() {
	if w.Min == 0 && w.Max == 0 {
		w.Max = 100
	}
}

// UpdateValue clamps the field Value to the range [Min,Max].
func (w *Slider) UpdateValue() {
	if w.Value < w.Min {
		w.Value = w.Min
	} else if w.Value > w.Max {
		w.Value = w.Max
	}
}

func (*sliderElement) Kind() *base.Kind {
	return &sliderKind
}

func (w *sliderElement) UpdateProps(data base.Widget) error {
	slider := data.(*Slider)

	// Fill in default values for the range.
	slider.UpdateRange()
	slider.UpdateValue()
	// Forward to the platform-dependant code
	return w.updateProps(slider)
}
