// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"clipster/goey/base"
	"clipster/goey/internal/gtk"
)

type textinputElement struct {
	Control

	onChange   func(string)
	onFocus    func()
	onBlur     func()
	onEnterKey func(string)
}

func (w *TextInput) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountTextbox(parent.Handle, w.Value, w.Placeholder, w.Disabled, w.Password, w.ReadOnly,
		w.OnChange != nil, w.OnFocus != nil, w.OnBlur != nil, w.OnEnterKey != nil)

	retval := &textinputElement{
		Control:    Control{control},
		onChange:   w.OnChange,
		onFocus:    w.OnFocus,
		onBlur:     w.OnBlur,
		onEnterKey: w.OnEnterKey,
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}

func (w *textinputElement) OnChange(value string) {
	if w.onChange != nil {
		w.onChange(value)
	}
}

func (w *textinputElement) OnFocus() {
	w.onFocus()
}

func (w *textinputElement) OnBlur() {
	w.onBlur()
}

func (w *textinputElement) OnEnterKey(value string) {
	w.onEnterKey(value)
}

func (w *textinputElement) Props() base.Widget {
	return &TextInput{
		Value:       gtk.TextboxText(w.handle),
		Placeholder: gtk.TextboxPlaceholder(w.handle),
		Disabled:    !gtk.WidgetSensitive(w.handle),
		Password:    gtk.TextboxPassword(w.handle),
		ReadOnly:    gtk.TextboxReadOnly(w.handle),
		OnChange:    w.onChange,
		OnFocus:     w.onFocus,
		OnBlur:      w.onBlur,
		OnEnterKey:  w.onEnterKey,
	}
}

func (w *textinputElement) updateProps(data *TextInput) error {
	w.onChange = nil // temporarily break OnChange to prevent event
	gtk.TextboxUpdate(w.handle, data.Value, data.Placeholder, data.Disabled,
		data.Password, data.ReadOnly,
		data.OnChange != nil, data.OnFocus != nil, data.OnBlur != nil, data.OnEnterKey != nil)

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	w.onEnterKey = data.OnEnterKey

	return nil
}
