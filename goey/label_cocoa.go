// +build cocoa darwin,!gtk

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/cocoa"
)

type labelElement struct {
	control *cocoa.Text
}

func (w *Label) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewText(parent.Handle, w.Text)

	retval := &labelElement{
		control: control,
	}
	return retval, nil
}

func (w *labelElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *labelElement) Layout(bc base.Constraints) base.Size {
	px := w.MinIntrinsicWidth(base.Inf)
	h := w.MinIntrinsicHeight(base.Inf)
	return bc.Constrain(base.Size{px, h})
}

func (w *labelElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 20 * base.DIP
}

func (w *labelElement) MinIntrinsicWidth(base.Length) base.Length {
	return 200 * base.DIP
}

func (w *labelElement) Props() base.Widget {
	return &Label{
		Text: w.control.Text(),
	}
}

func (w *labelElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *labelElement) updateProps(data *Label) error {
	w.control.SetText(data.Text)
	return nil
}
