package goey

import (
	"guitest/goey/base"
)

var (
	paddingKind = base.NewKind("guitest/goey.Padding")
)

// Insets describe padding that should be added around a widget.
type Insets struct {
	Top    base.Length
	Right  base.Length
	Bottom base.Length
	Left   base.Length
}

// DefaultInsets returns the (perhaps platform-dependent) default insets for
// widgets inside of a top-level window.
func DefaultInsets() Insets {
	const padding = 11 * base.DIP
	return Insets{padding, padding, padding, padding}
}

// UniformInsets returns a padding description where the padding is equal on
// all four sides.
func UniformInsets(l base.Length) Insets {
	return Insets{l, l, l, l}
}

// String returns a string representation of i like "Insets{92:00,92:00,46:00,46:00}".
func (i *Insets) String() string {
	return "Insets{" + i.Top.String() + "," + i.Right.String() + "," +
		i.Bottom.String() + "," + i.Left.String() + "}"
}

// Padding describes a widget that adds some space around a single child widget.
//
// The size of the control will match the size of the child element, although
// padding will be added between the border of the padding and the child
// element as specified by the field Insets.
type Padding struct {
	Insets Insets      // Space between edge of element and the child element.
	Child  base.Widget // Child widget.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Padding) Kind() *base.Kind {
	return &paddingKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Padding) Mount(parent base.Control) (base.Element, error) {
	child, err := base.Mount(parent, w.Child)
	if err != nil {
		return nil, err
	}

	return &paddingElement{
		parent: parent,
		child:  child,
		insets: w.Insets,
	}, nil
}

type paddingElement struct {
	parent    base.Control
	child     base.Element
	childSize base.Size
	insets    Insets
}

func (w *paddingElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
}

func (*paddingElement) Kind() *base.Kind {
	return &paddingKind
}

func (w *paddingElement) Layout(bc base.Constraints) base.Size {
	hinset := w.insets.Left + w.insets.Right
	vinset := w.insets.Top + w.insets.Bottom

	innerConstraints := bc.Inset(hinset, vinset)
	w.childSize = w.child.Layout(innerConstraints)
	return base.Size{
		w.childSize.Width + hinset,
		w.childSize.Height + vinset,
	}
}

func (w *paddingElement) MinIntrinsicHeight(width base.Length) base.Length {
	vinset := w.insets.Top + w.insets.Bottom
	return w.child.MinIntrinsicHeight(width) + vinset
}

func (w *paddingElement) MinIntrinsicWidth(height base.Length) base.Length {
	hinset := w.insets.Left + w.insets.Right
	return w.child.MinIntrinsicWidth(height) + hinset
}

func (w *paddingElement) SetBounds(bounds base.Rectangle) {
	bounds.Min.X += w.insets.Left
	bounds.Min.Y += w.insets.Top
	bounds.Max.X -= w.insets.Right
	bounds.Max.Y -= w.insets.Bottom
	w.child.SetBounds(bounds)
}

func (w *paddingElement) updateProps(data *Padding) (err error) {
	w.child, err = base.DiffChild(w.parent, w.child, data.Child)
	w.insets = data.Insets
	return err
}

func (w *paddingElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*Padding))
}
