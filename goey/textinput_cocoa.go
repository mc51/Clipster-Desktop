// +build cocoa darwin,!gtk

package goey

import (
	"time"

	"guitest/goey/base"
	"guitest/goey/internal/cocoa"
	"guitest/goey/loop"
)

type textinputElement struct {
	control *cocoa.TextField
}

func (w *TextInput) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewTextField(parent.Handle, w.Value, w.Password)
	control.SetPlaceholder(w.Placeholder)
	control.SetEnabled(!w.Disabled)
	control.SetEditable(!w.ReadOnly)
	control.SetCallbacks(w.OnChange, w.OnFocus, w.OnBlur, w.OnEnterKey)

	retval := &textinputElement{
		control: control,
	}
	return retval, nil
}

func (w *textinputElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *textinputElement) Layout(bc base.Constraints) base.Size {
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *textinputElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *textinputElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *textinputElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *textinputElement) TakeFocus() bool {
	return w.control.MakeFirstResponder()
}

// TypeKeys sends events to the control as if the string was typed by a user.
func (w *textinputElement) TypeKeys(text string) chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		time.Sleep(500 * time.Millisecond)
		for _, r := range text {
			err := loop.Do(func() error {
				w.control.SendKey(uint(r))
				return nil
			})
			if err != nil {
				errs <- err
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	return errs
}

func (w *textinputElement) Props() base.Widget {
	onchange, onfocus, onblur, onenterkey := w.control.Callbacks()

	return &TextInput{
		Value:       w.control.Value(),
		Disabled:    !w.control.IsEnabled(),
		Placeholder: w.control.Placeholder(),
		Password:    w.control.IsPassword(),
		ReadOnly:    !w.control.IsEditable(),
		OnChange:    onchange,
		OnFocus:     onfocus,
		OnBlur:      onblur,
		OnEnterKey:  onenterkey,
	}
}

func (w *textinputElement) updateProps(data *TextInput) error {
	w.control.SetValue(data.Value)
	w.control.SetPlaceholder(data.Placeholder)
	w.control.SetEnabled(!data.Disabled)
	w.control.SetEditable(!data.ReadOnly)
	w.control.SetCallbacks(data.OnChange, data.OnFocus, data.OnBlur, data.OnEnterKey)
	return nil
}
