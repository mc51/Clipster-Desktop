package goey

import (
	"clipster/goey/base"
)

var (
	hrKind = base.NewKind("clipster/goey.HR")
)

// HR describes a widget that is a horizontal separator.
type HR struct {
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*HR) Kind() *base.Kind {
	return &hrKind
}

// Mount creates a horizontal rule control in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *HR) Mount(parent base.Control) (base.Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*hrElement) Kind() *base.Kind {
	return &hrKind
}

func (w *hrElement) Props() base.Widget {
	return &HR{}
}

func (*hrElement) Layout(bc base.Constraints) base.Size {
	if bc.HasBoundedWidth() {
		return bc.Constrain(base.Size{bc.Max.Width, 13 * DIP})
	}
	return bc.Constrain(base.Size{13 * DIP, 13 * DIP})
}

func (w *hrElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 13 * DIP
}

func (w *hrElement) MinIntrinsicWidth(height base.Length) base.Length {
	return 13 * DIP
}

func (w *hrElement) UpdateProps(data base.Widget) error {
	// This widget does not have any properties, so there cannot be anything
	// to update.
	return nil
}
