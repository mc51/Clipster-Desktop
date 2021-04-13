package goey

import (
	"image/color"

	"guitest/goey/base"
)

var (
	decorationKind = base.NewKind("guitest/goey.Decoration")
)

// Decoration describes a widget that provides a border and background, and
// possibly containing a single child widget.
//
// The size of the control will match the size of the child element, although
// padding will be added between the border of the decoration and the child
// element as specified by the field Insets.
type Decoration struct {
	Fill   color.RGBA  // Background color used to fill interior.
	Stroke color.RGBA  // Stroke color used to draw outline.
	Insets Insets      // Space between border of the decoration and the child element.
	Radius base.Length // Radius of the widgets corners.
	Child  base.Widget // Child widget.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Decoration) Kind() *base.Kind {
	return &decorationKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Decoration) Mount(parent base.Control) (base.Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*decorationElement) Kind() *base.Kind {
	return &decorationKind
}

func (w *decorationElement) Layout(bc base.Constraints) base.Size {
	hinset := w.insets.Left + w.insets.Right
	vinset := w.insets.Top + w.insets.Bottom

	innerConstraints := bc.Inset(hinset, vinset)
	w.childSize = w.child.Layout(innerConstraints)
	return base.Size{
		max(w.childSize.Width+hinset, base.FromPixelsX(2)),
		max(w.childSize.Height+vinset, base.FromPixelsY(2)),
	}
}

func (w *decorationElement) MinIntrinsicHeight(width base.Length) base.Length {
	vinset := w.insets.Top + w.insets.Bottom
	return max(w.child.MinIntrinsicHeight(width)+vinset,
		base.FromPixelsX(2))
}

func (w *decorationElement) MinIntrinsicWidth(height base.Length) base.Length {
	hinset := w.insets.Left + w.insets.Right
	return max(w.child.MinIntrinsicWidth(height)+hinset,
		base.FromPixelsY(2))

}

func (w *decorationElement) UpdateProps(data base.Widget) error {
	// Forward to the platform-dependant code
	return w.updateProps(data.(*Decoration))
}
