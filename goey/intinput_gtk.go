// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/gtk"
)

type intinputElement struct {
	Control

	onChange   func(int64)
	onFocus    func()
	onBlur     func()
	onEnterKey func(int64)
}

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountIntInput(parent.Handle, w.Value, w.Placeholder, w.Disabled,
		w.Min, w.Max,
		w.OnChange != nil, w.OnFocus != nil, w.OnBlur != nil, w.OnEnterKey != nil)

	// Create the element
	retval := &intinputElement{
		Control:    Control{control},
		onChange:   w.OnChange,
		onFocus:    w.OnFocus,
		onBlur:     w.OnBlur,
		onEnterKey: w.OnEnterKey,
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}

// Because GTK uses double (or float64 in Go) to store the range, it cannot
// keep full precision for values near the minimum or maximum of the int64
// range.  This function is just for Props, and is used to adjust the int64
// to get a match.
func toInt64(scale float64) int64 {
	a := int64(scale)
	if float64(a) == scale {
		return a
	}
	if float64(a-1) == scale {
		return a - 1
	}
	println("mismatch", a, scale)
	return a
}

func (w *intinputElement) OnChange(value int64) {
	if w.onChange != nil {
		w.onChange(value)
	}
}

func (w *intinputElement) OnFocus() {
	w.onFocus()
}

func (w *intinputElement) OnBlur() {
	w.onBlur()
}

func (w *intinputElement) OnEnterKey(value int64) {
	w.onEnterKey(value)
}

func (w *intinputElement) Props() base.Widget {
	return &IntInput{
		Value:       gtk.IntinputValue(w.handle),
		Placeholder: gtk.TextboxPlaceholder(w.handle),
		Disabled:    !gtk.WidgetSensitive(w.handle),
		Min:         toInt64(gtk.IntinputMin(w.handle)),
		Max:         toInt64(gtk.IntinputMax(w.handle)),
		OnChange:    w.onChange,
		OnFocus:     w.onFocus,
		OnBlur:      w.onBlur,
		OnEnterKey:  w.onEnterKey,
	}
}

func (w *intinputElement) updateProps(data *IntInput) error {
	gtk.IntinputUpdate(w.handle, data.Value, data.Placeholder, data.Disabled,
		data.Min, data.Max,
		data.OnChange != nil, data.OnFocus != nil, data.OnBlur != nil, data.OnEnterKey != nil)

	return nil
}
