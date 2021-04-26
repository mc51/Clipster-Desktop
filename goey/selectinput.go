package goey

import (
	"clipster/goey/base"
)

var (
	selectKind = base.NewKind("clipster/goey.SelectInput")
)

// SelectInput describes a widget that users can click to select one from a fixed list of choices.
type SelectInput struct {
	Items    []string        // Items is an array of strings representing the user's possible choices
	Value    int             // Value is the index of the currently selected item
	Unset    bool            // Unset is a flag indicating that no choice has yet been made
	Disabled bool            // Disabled is a flag indicating that the user cannot interact with this field
	OnChange func(value int) // OnChange will be called whenever the user changes the value for this field
	OnFocus  func()          // OnFocus will be called whenever the field receives the keyboard focus
	OnBlur   func()          // OnBlur will be called whenever the field loses the keyboard focus
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*SelectInput) Kind() *base.Kind {
	return &selectKind
}

// Mount creates a select control (combobox) in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *SelectInput) Mount(parent base.Control) (base.Element, error) {
	// Update Value and Unset to make sure that are they are coherent with the
	// length of items.
	w.UpdateValue()

	// Forward to the platform-dependant code
	return w.mount(parent)
}

// UpdateValue will ensure that the Value is within the range of choices
// provided by w.Items.
func (w *SelectInput) UpdateValue() {
	if length := len(w.Items); length > 0 {
		if w.Value >= length {
			w.Value = length - 1
		} else if w.Value < 0 {
			w.Value = 0
		}
	} else {
		w.Value = 0
		w.Unset = true
	}
}

func (*selectinputElement) Kind() *base.Kind {
	return &selectKind
}

func (w *selectinputElement) UpdateProps(data base.Widget) error {
	si := data.(*SelectInput)

	// Update Value and Unset to make sure that are they are coherent.
	si.UpdateValue()
	// Forward to the platform-dependant code
	return w.updateProps(si)
}
