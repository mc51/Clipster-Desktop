package goey

import (
	"guitest/goey/base"
)

var (
	expandKind = base.NewKind("guitest/goey.Expand")
)

// Expand wraps another widget to indicate that the widget should expand to
// occupy any available space in a HBox or VBox.  When used in any other
// context, the widget will be ignored, and behavior delegated to the child
// widget.
//
// In an HBox or VBox, the widget will be positioned according to the rules
// of its child.  However, any excess space along the main axis will be added
// based on the ratio of the widget's factor to the sum of factors for all
// widgets in the box.  Note that
type Expand struct {
	Factor int         // Fraction (minus one) of available space used by this widget
	Child  base.Widget // Child widget.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Expand) Kind() *base.Kind {
	return &expandKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Expand) Mount(parent base.Control) (base.Element, error) {
	// Mount the child
	child, err := base.Mount(parent, w.Child)
	if err != nil {
		return nil, err
	}

	return &expandElement{
		parent: parent,
		child:  child,
		factor: w.Factor,
	}, nil
}

type expandElement struct {
	parent base.Control
	child  base.Element
	factor int
}

func (w *expandElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
}

func (*expandElement) Kind() *base.Kind {
	return &expandKind
}

func (w *expandElement) Layout(bc base.Constraints) base.Size {
	return w.child.Layout(bc)
}

func (w *expandElement) MinIntrinsicHeight(width base.Length) base.Length {
	return w.child.MinIntrinsicHeight(width)
}

func (w *expandElement) MinIntrinsicWidth(height base.Length) base.Length {
	return w.child.MinIntrinsicWidth(height)
}

func (w *expandElement) SetBounds(bounds base.Rectangle) {
	w.child.SetBounds(bounds)
}

func (w *expandElement) updateProps(data *Expand) (err error) {
	w.child, err = base.DiffChild(w.parent, w.child, data.Child)
	w.factor = data.Factor
	return err
}

func (w *expandElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*Expand))
}
