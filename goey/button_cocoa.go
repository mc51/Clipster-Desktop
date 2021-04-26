// +build cocoa darwin,!gtk

package goey

import (
	"clipster/goey/base"
	"clipster/goey/internal/cocoa"
)

type buttonElement struct {
	control *cocoa.Button
}

func (w *Button) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewButton(parent.Handle, w.Text)
	control.SetCallbacks(w.OnClick, nil, w.OnFocus, w.OnBlur)
	control.SetEnabled(!w.Disabled)
	control.SetDefault(w.Default)

	retval := &buttonElement{
		control: control,
	}
	return retval, nil
}

func (w *buttonElement) Click() {
	w.control.PerformClick()
}

func (w *buttonElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *buttonElement) Layout(bc base.Constraints) base.Size {
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *buttonElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *buttonElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *buttonElement) Props() base.Widget {
	onclick, _, onfocus, onblur := w.control.Callbacks()

	return &Button{
		Text:     w.control.Title(),
		Disabled: !w.control.IsEnabled(),
		Default:  w.control.IsDefault(),
		OnClick:  onclick,
		OnFocus:  onfocus,
		OnBlur:   onblur,
	}
}

func (w *buttonElement) TakeFocus() bool {
	return w.control.MakeFirstResponder()
}

func (w *buttonElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *buttonElement) updateProps(data *Button) error {
	w.control.SetTitle(data.Text)
	w.control.SetEnabled(!data.Disabled)
	w.control.SetDefault(data.Default)
	w.control.SetCallbacks(data.OnClick, nil, data.OnFocus, data.OnBlur)
	return nil
}
