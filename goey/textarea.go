package goey

import (
	"clipster/goey/base"
)

var (
	textareaKind = base.NewKind("clipster/goey.TextArea")
)

// TextArea describes a widget that users input or update a multi-line of text.
// The model for the value is a string value.
//
// Using a placeholder may not be supported on all platforms.  No errors will
// be generated, but the placeholder text may not appear on screen.
type TextArea struct {
	Value       string             // Values is the current string for the field
	Placeholder string             // Placeholder is a descriptive text that can be displayed when the field is empty
	Disabled    bool               // Disabled is a flag indicating that the user cannot interact with this field
	ReadOnly    bool               // ReadOnly is a flag indicate that the contents cannot be modified by the user
	MinLines    int                // MinLines describes the minimum number of lines that should be visible for layout
	OnChange    func(value string) // OnChange will be called whenever the user changes the value for this field
	OnFocus     func()             // OnFocus will be called whenever the field receives the keyboard focus
	OnBlur      func()             // OnBlur will be called whenever the field loses the keyboard focus
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*TextArea) Kind() *base.Kind {
	return &textareaKind
}

// Mount creates a text area control in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *TextArea) Mount(parent base.Control) (base.Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*textareaElement) Kind() *base.Kind {
	return &textareaKind
}

func (w *textareaElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*TextArea))
}

func minlinesDefault(value int) int {
	if value < 1 {
		return 3
	}
	return value
}
