package goey

import (
	"guitest/goey/base"
)

var (
	alignKind = base.NewKind("guitest/goey.Align")
)

// Alignment represents the position of a child widget along one dimension.
// Some common values for alignment, such as AlignStart, AlignCenter, and AlignEnd,
// are given constants, but other values are possible.  For example, to align
// a child with an position of 25%, use (AlignStart + AlignCenter) / 2.
type Alignment int16

// Common values for alignment, representing the position of child widget.
const (
	AlignStart  Alignment = -32768 // Widget is aligned at the start (left or top).
	AlignCenter Alignment = 0      // Widget is aligned at the center.
	AlignEnd    Alignment = 0x7fff // Widget is aligned at the end (right or bottom).
)

// Align describes a widget that aligns a single child widget within its borders.
//
// The default position is for the child widget to be centered.  To change the
// position of the child, the horizontal and vertical alignment (the fields
// HAlign and VAlign) should be adjusted.
//
// The size of the control depends on the WidthFactor and HeightFactor.  If zero,
// the widget will try to be as large as possible or match the child, depending
// on whether the box constraints are bound or not.  If a factor is greater than
// zero, then the widget will try to size itself to be that much larger than the
// child widget.
type Align struct {
	HAlign       Alignment   // Horizontal alignment of child widget.
	VAlign       Alignment   // Vertical alignment of child widget.
	WidthFactor  float64     // If greater than zero, ratio of container width to child width.
	HeightFactor float64     // If greater than zero, ratio of container height to child height.
	Child        base.Widget // Child widget.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Align) Kind() *base.Kind {
	return &alignKind
}

// Mount creates an aligned layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Align) Mount(parent base.Control) (base.Element, error) {
	// Mount the child
	child, err := base.Mount(parent, w.Child)
	if err != nil {
		return nil, err
	}

	return &alignElement{
		parent:       parent,
		child:        child,
		widthFactor:  w.WidthFactor,
		heightFactor: w.HeightFactor,
		hAlign:       w.HAlign,
		vAlign:       w.VAlign,
	}, nil
}

type alignElement struct {
	parent       base.Control
	child        base.Element
	childSize    base.Size
	hAlign       Alignment
	vAlign       Alignment
	widthFactor  float64
	heightFactor float64
}

func (w *alignElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
}

func (*alignElement) Kind() *base.Kind {
	return &alignKind
}

func (w *alignElement) Layout(bc base.Constraints) base.Size {
	size := w.child.Layout(bc.Loosen())
	w.childSize = size
	if w.widthFactor > 0 {
		size.Width = base.Length(float64(size.Width) * w.widthFactor)
	}
	if w.heightFactor > 0 {
		size.Height = base.Length(float64(size.Height) * w.heightFactor)
	}
	return bc.Constrain(size)
}

func (w *alignElement) MinIntrinsicHeight(width base.Length) base.Length {
	height := w.child.MinIntrinsicHeight(width)
	if w.heightFactor > 0 {
		return base.Length(float64(height) * w.heightFactor)
	}
	return height
}

func (w *alignElement) MinIntrinsicWidth(height base.Length) base.Length {
	width := w.child.MinIntrinsicWidth(height)
	if w.widthFactor > 0 {
		return base.Length(float64(width) * w.widthFactor)
	}
	return width
}

func (w *alignElement) SetBounds(bounds base.Rectangle) {
	x := bounds.Min.X.Scale(int(w.hAlign)-int(AlignEnd), int(AlignStart)-int(AlignEnd)) +
		(bounds.Max.X-w.childSize.Width).Scale(int(w.hAlign)-int(AlignStart), int(AlignEnd)-int(AlignStart))
	y := bounds.Min.Y.Scale(int(w.vAlign)-int(AlignEnd), int(AlignStart)-int(AlignEnd)) +
		(bounds.Max.Y-w.childSize.Height).Scale(int(w.vAlign)-int(AlignStart), int(AlignEnd)-int(AlignStart))
	w.child.SetBounds(base.Rectangle{
		base.Point{x, y},
		base.Point{x + w.childSize.Width, y + w.childSize.Height},
	})
}

func (w *alignElement) updateProps(data *Align) (err error) {
	w.child, err = base.DiffChild(w.parent, w.child, data.Child)
	w.widthFactor = data.WidthFactor
	w.heightFactor = data.HeightFactor
	w.hAlign = data.HAlign
	w.vAlign = data.VAlign
	return err
}

func (w *alignElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*Align))
}
