// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/gtk"
)

type buttonElement struct {
	Control

	onClick func()
	onFocus func()
	onBlur  func()
}

func (w *Button) mount(parent base.Control) (base.Element, error) {
	// Create the control
	handle := gtk.MountButton(parent.Handle, w.Text, w.Disabled, w.Default,
		w.OnClick != nil, w.OnFocus != nil, w.OnBlur != nil)

	// Create the element
	retval := &buttonElement{
		Control: Control{handle},
		onClick: w.OnClick,
		onFocus: w.OnFocus,
		onBlur:  w.OnBlur,
	}
	gtk.RegisterWidget(handle, retval)

	return retval, nil
}

func (w *buttonElement) Click() {
	gtk.ButtonClick(w.handle)
}

func (w *buttonElement) OnClick() {
	w.onClick()
}

func (w *buttonElement) OnFocus() {
	w.onFocus()
}

func (w *buttonElement) OnBlur() {
	w.onBlur()
}

func (w *buttonElement) Props() base.Widget {
	return &Button{
		Text:     gtk.ButtonText(w.handle),
		Disabled: !gtk.WidgetSensitive(w.handle),
		Default:  gtk.WidgetCanDefault(w.handle),
		OnClick:  w.onClick,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *buttonElement) updateProps(data *Button) error {
	gtk.ButtonUpdate(w.handle, data.Text, data.Disabled, data.Default,
		data.OnClick != nil, data.OnFocus != nil, data.OnBlur != nil)

	w.onClick = data.OnClick
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}
