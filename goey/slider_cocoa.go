// +build cocoa darwin,!gtk

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/cocoa"
)

type sliderElement struct {
	control *cocoa.Slider
}

func (w *Slider) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewSlider(parent.Handle, w.Min, w.Value, w.Max)
	control.SetEnabled(!w.Disabled)
	control.SetCallbacks(w.OnChange, w.OnFocus, w.OnBlur)

	retval := &sliderElement{
		control: control,
	}
	return retval, nil
}

func (w *sliderElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *sliderElement) Layout(bc base.Constraints) base.Size {
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *sliderElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *sliderElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *sliderElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *sliderElement) Props() base.Widget {
	onchange, onfocus, onblur := w.control.Callbacks()
	return &Slider{
		Value:    w.control.Value(),
		Disabled: !w.control.IsEnabled(),
		Min:      w.control.Min(),
		Max:      w.control.Max(),
		OnChange: onchange,
		OnFocus:  onfocus,
		OnBlur:   onblur,
	}
}

func (w *sliderElement) TakeFocus() bool {
	return w.control.MakeFirstResponder()
}

func (w *sliderElement) updateProps(data *Slider) error {
	w.control.Update(data.Min, data.Value, data.Max)
	w.control.SetEnabled(!data.Disabled)
	return nil
}
