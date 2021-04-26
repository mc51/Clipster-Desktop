// +build cocoa darwin,!gtk

package goey

import (
	"clipster/goey/base"
	"clipster/goey/internal/cocoa"
)

type checkboxElement struct {
	control *cocoa.Button
}

func (w *Checkbox) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewCheckButton(parent.Handle, w.Text, w.Value)
	control.SetCallbacks(nil, w.OnChange, w.OnFocus, w.OnBlur)
	control.SetEnabled(!w.Disabled)

	retval := &checkboxElement{
		control: control,
	}
	return retval, nil
}

func (w *checkboxElement) Click() {
	w.control.PerformClick()
}

func (w *checkboxElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *checkboxElement) Layout(bc base.Constraints) base.Size {
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *checkboxElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *checkboxElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *checkboxElement) Props() base.Widget {
	text := w.control.Title()
	value := w.control.State()
	disabled := !w.control.IsEnabled()
	_, onchange, onfocus, onblur := w.control.Callbacks()

	return &Checkbox{
		Text:     text,
		Value:    value,
		Disabled: disabled,
		OnChange: onchange,
		OnFocus:  onfocus,
		OnBlur:   onblur,
	}
}

func (w *checkboxElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *checkboxElement) TakeFocus() bool {
	return w.control.MakeFirstResponder()
}

func (w *checkboxElement) updateProps(data *Checkbox) error {
	w.control.SetTitle(data.Text)
	w.control.SetState(data.Value)
	w.control.SetEnabled(!data.Disabled)
	w.control.SetCallbacks(nil, data.OnChange, data.OnFocus, data.OnBlur)
	return nil
}
