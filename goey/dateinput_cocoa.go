// +build cocoa darwin,!gtk

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/cocoa"
)

type dateinputElement struct {
	control *cocoa.Text
}

func (w *DateInput) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewText(parent.Handle, "date input")

	retval := &dateinputElement{
		control: control,
	}
	return retval, nil
}

func (w *dateinputElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *dateinputElement) Layout(bc base.Constraints) base.Size {
	px := w.MinIntrinsicWidth(base.Inf)
	h := w.MinIntrinsicHeight(base.Inf)
	return bc.Constrain(base.Size{px, h})
}

func (w *dateinputElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 20 * base.DIP
}

func (w *dateinputElement) MinIntrinsicWidth(base.Length) base.Length {
	return 200 * base.DIP
}

func (w *dateinputElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *dateinputElement) updateProps(data *DateInput) error {
	return nil
}
