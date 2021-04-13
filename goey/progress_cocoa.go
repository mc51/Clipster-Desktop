// +build cocoa darwin,!gtk

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/cocoa"
)

type progressElement struct {
	control *cocoa.Progress
}

func (w *Progress) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewProgress(parent.Handle,
		float64(w.Min), float64(w.Value), float64(w.Max))

	retval := &progressElement{
		control: control,
	}
	return retval, nil
}

func (w *progressElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *progressElement) Layout(bc base.Constraints) base.Size {
	px := w.MinIntrinsicWidth(base.Inf)
	h := w.MinIntrinsicHeight(base.Inf)
	return bc.Constrain(base.Size{px, h})
}

func (w *progressElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 20 * base.DIP
}

func (w *progressElement) MinIntrinsicWidth(base.Length) base.Length {
	return 200 * base.DIP
}

func (w *progressElement) Props() base.Widget {
	min := w.control.Min()
	value := w.control.Value()
	max := w.control.Max()
	return &Progress{
		Value: int(value),
		Min:   int(min),
		Max:   int(max),
	}
}

func (w *progressElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *progressElement) updateProps(data *Progress) error {
	w.control.Update(float64(data.Min), float64(data.Value), float64(data.Max))
	return nil
}
