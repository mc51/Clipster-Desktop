// +build cocoa darwin,!gtk

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/cocoa"
)

type paragraphElement struct {
	control *cocoa.Text
}

func (w *P) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewText(parent.Handle, w.Text)
	control.SetAlignment(int(w.Align))

	retval := &paragraphElement{
		control: control,
	}
	return retval, nil
}

func (w *paragraphElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *paragraphElement) measureReflowLimits() {
	x := w.control.EightyEms()
	paragraphMaxWidth = base.FromPixelsX(x)
}

func (w *paragraphElement) MinIntrinsicHeight(width base.Length) base.Length {
	if width == base.Inf {
		width = w.maxReflowWidth()
	}

	y := w.control.MinHeight(width.PixelsX())
	return base.FromPixelsY(y)
}

func (w *paragraphElement) MinIntrinsicWidth(height base.Length) base.Length {
	if height != base.Inf {
		panic("not implemented")
	}

	x := w.control.MinWidth()
	return min(base.FromPixelsX(x), w.minReflowWidth())
}

func (w *paragraphElement) Props() base.Widget {
	return &P{
		Text:  w.control.Text(),
		Align: TextAlignment(w.control.Alignment()),
	}
}

func (w *paragraphElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *paragraphElement) updateProps(data *P) error {
	w.control.SetText(data.Text)
	w.control.SetAlignment(int(data.Align))
	return nil
}
