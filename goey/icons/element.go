package icons

import (
	"clipster/goey"
	"clipster/goey/base"
)

type iconElement struct {
	parent base.Control
	child  base.Element
	icon   rune
}

func (w *iconElement) Close() {
	w.child.Close()
	w.child = nil
}

func (*iconElement) Kind() *base.Kind {
	return &kind
}

func (w *iconElement) Layout(bc base.Constraints) base.Size {
	return w.child.Layout(bc)
}

func (w *iconElement) MinIntrinsicHeight(width base.Length) base.Length {
	return w.child.MinIntrinsicHeight(width)
}

func (w *iconElement) MinIntrinsicWidth(height base.Length) base.Length {
	return w.child.MinIntrinsicWidth(height)
}

func (w *iconElement) SetBounds(bounds base.Rectangle) {
	w.child.SetBounds(bounds)
}

func (w *iconElement) updateProps(data Icon) error {
	if rune(data) == w.icon {
		return nil
	}

	// TODO:  We should either cache and reuse the image data, or at least
	// draw onto the existing buffer.
	img, err := DrawImage(rune(data))
	if err != nil {
		return err
	}

	elem, err := base.DiffChild(w.parent, w.child, &goey.Img{Image: img})
	if err != nil {
		return err
	}
	w.child = elem
	w.icon = rune(data)
	return nil
}

func (w *iconElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(Icon))
}
