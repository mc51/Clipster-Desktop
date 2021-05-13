// +build gtk linux darwin freebsd openbsd

package goey

import (
	"clipster/goey/base"
	"clipster/goey/internal/gtk"
)

type checkboxElement struct {
	Control

	onChange func(bool)
	onFocus  func()
	onBlur   func()
}

func (w *Checkbox) mount(parent base.Control) (base.Element, error) {
	// Create the control
	control := gtk.MountCheckbox(parent.Handle, w.Value, w.Text, w.Disabled,
		w.OnChange != nil, w.OnFocus != nil, w.OnBlur != nil)

	// Create the element
	retval := &checkboxElement{
		Control:  Control{control},
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}

func (w *checkboxElement) Click() {
	gtk.CheckboxClick(w.handle)
}

func (w *checkboxElement) OnChange(value bool) {
	if w.onChange != nil {
		w.onChange(value)
	}
}

func (w *checkboxElement) OnFocus() {
	w.onFocus()
}

func (w *checkboxElement) OnBlur() {
	w.onBlur()
}

func (w *checkboxElement) Props() base.Widget {
	return &Checkbox{
		Value:    gtk.CheckboxValue(w.handle),
		Text:     gtk.ButtonText(w.handle),
		Disabled: !gtk.WidgetSensitive(w.handle),
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *checkboxElement) updateProps(data *Checkbox) error {
	w.onChange = nil // temporarily break OnChange to prevent event
	gtk.CheckboxUpdate(w.handle, data.Value, data.Text, data.Disabled,
		data.OnChange != nil, data.OnFocus != nil, data.OnBlur != nil)

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}
