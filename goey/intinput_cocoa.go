// +build cocoa darwin,!gtk

package goey

import (
	"clipster/goey/base"
	"clipster/goey/internal/cocoa"
)

type intinputElement struct {
	control *cocoa.IntField
}

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewIntField(parent.Handle, w.Value, w.Min, w.Max)
	control.SetPlaceholder(w.Placeholder)
	control.SetEnabled(!w.Disabled)
	control.SetCallbacks(w.OnChange, w.OnFocus, w.OnBlur)

	retval := &intinputElement{
		control: control,
	}
	return retval, nil
}

func (w *intinputElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *intinputElement) Layout(bc base.Constraints) base.Size {
	px := w.MinIntrinsicWidth(base.Inf)
	h := w.MinIntrinsicHeight(base.Inf)
	return bc.Constrain(base.Size{px, h})
}

func (w *intinputElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *intinputElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *intinputElement) Props() base.Widget {
	onchange, onfocus, onblur := w.control.Callbacks()

	return &IntInput{
		Value:       w.control.Value(),
		Min:         w.control.Min(),
		Max:         w.control.Max(),
		Disabled:    !w.control.IsEnabled(),
		Placeholder: w.control.Placeholder(),
		OnChange:    onchange,
		OnFocus:     onfocus,
		OnBlur:      onblur,
	}
}

func (w *intinputElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *intinputElement) updateProps(data *IntInput) error {
	w.control.SetValue(data.Value, data.Min, data.Max)
	w.control.SetPlaceholder(data.Placeholder)
	w.control.SetEnabled(!data.Disabled)
	w.control.SetCallbacks(data.OnChange, data.OnFocus, data.OnBlur)
	return nil
}
