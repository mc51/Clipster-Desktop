package mock

import (
	"guitest/goey/base"
)

var (
	mockKind = base.NewKind("guitest/goey/mock.Widget")
)

// New returns a new mock element.
func New(size base.Size) *Element {
	return &Element{
		Size: size,
	}
}

// NewIfNotZero returns a new mock element if the size is not zero, otherwise
// it returns the nil element.
func NewIfNotZero(size base.Size) base.Element {
	if size.IsZero() {
		elem, err := base.Mount(base.Control{}, nil)
		if err != nil {
			panic("Mounting a memory only element should never fail")
		}
		return elem
	}

	return New(size)
}

// NewList returns a slice of new mock elements.
func NewList(sizes ...base.Size) []base.Element {
	ret := make([]base.Element, 0, len(sizes))
	for _, v := range sizes {
		ret = append(ret, &Element{Size: v})
	}
	return ret
}

// Widget is a mock widget.  When mounted, it will create a mock element.
type Widget struct {
	Size base.Size
	Err  error
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Widget) Kind() *base.Kind {
	return &mockKind
}

// Mount creates an mock control.
func (w *Widget) Mount(parent base.Control) (base.Element, error) {
	// Check if the widget is supposed to fail when mounted.
	if w.Err != nil {
		return nil, w.Err
	}

	// Create a mock element.
	return &Element{
		Size: w.Size,
	}, nil
}

// Element is a mock element.
// Although this is a leaf element (i.e. it does not have any children),
// there is no control associated with this element.
type Element struct {
	Size base.Size

	bounds base.Rectangle
	closed bool
}

// Close removes the widget from the GUI, and frees any associated resources.
// This is a no-op for a mock element.
func (w *Element) Close() {
	if w.closed {
		panic("Element already closed")
	}
	w.closed = true
}

// Kind returns the concrete type for the Element.
func (*Element) Kind() *base.Kind {
	return &mockKind
}

// Layout determines the best size for an element that satisfies the
// constraints.  For a mock element, it will try to match the specific
// size in the field Size.
func (w *Element) Layout(bc base.Constraints) base.Size {
	return bc.Constrain(w.Size)
}

// MinIntrinsicHeight returns the minimum height that this element requires
// to be correctly displayed.
func (w *Element) MinIntrinsicHeight(base.Length) base.Length {
	return w.Size.Height
}

// MinIntrinsicWidth returns the minimum width that this element requires
// to be correctly displayed.
func (w *Element) MinIntrinsicWidth(base.Length) base.Length {
	return w.Size.Width
}

// Props recreates the widget description for this element.
func (w *Element) Props() base.Widget {
	return &Widget{
		Size: w.Size,
	}
}

// Bounds returns the position of the widget.
func (w *Element) Bounds() base.Rectangle {
	return w.bounds
}

// SetBounds updates the position of the widget.
func (w *Element) SetBounds(bounds base.Rectangle) {
	w.bounds = bounds
}

func (w *Element) updateProps(data *Widget) error {
	w.Size = data.Size
	return nil
}

// UpdateProps will update the properties of the widget.
func (w *Element) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*Widget))
}
