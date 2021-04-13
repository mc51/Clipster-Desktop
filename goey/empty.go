package goey

import (
	"guitest/goey/base"
)

var (
	emptyKind = base.NewKind("guitest/goey.Empty")
)

// Empty describes a widget that is either a horizontal or vertical gap.
//
// The size of the control will be a (perhaps platform dependent) spacing
// between controls.  This applies to both the width and height.
type Empty struct {
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Empty) Kind() *base.Kind {
	return &emptyKind
}

// Mount creates a horizontal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Empty) Mount(parent base.Control) (base.Element, error) {
	retval := &emptyElement{}

	return retval, nil
}

type emptyElement struct {
}

func (w *emptyElement) Close() {
	// Virtual control, so no resources to release
}

func (*emptyElement) Kind() *base.Kind {
	return &emptyKind
}

func (w *emptyElement) Props() base.Widget {
	return &Empty{}
}

func (w *emptyElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *emptyElement) MinIntrinsicHeight(width base.Length) base.Length {
	// Same as static text
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13 * base.DIP
}

func (w *emptyElement) MinIntrinsicWidth(height base.Length) base.Length {
	// Same as static text
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13 * base.DIP
}

func (w *emptyElement) SetBounds(bounds base.Rectangle) {
	// Virtual control, so no resource to resize
}

func (w *emptyElement) UpdateProps(data base.Widget) error {
	// This widget does not have any properties, so there cannot be anything
	// to update.
	return nil
}
