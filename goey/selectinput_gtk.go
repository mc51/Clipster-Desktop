// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"bytes"

	"guitest/goey/base"
	"guitest/goey/internal/gtk"
)

type selectinputElement struct {
	Control

	onChange func(int)
	onFocus  func()
	onBlur   func()
}

func (w *SelectInput) serializeItems() string {
	buffer := bytes.Buffer{}

	for _, v := range w.Items {
		buffer.WriteString(v)
		buffer.WriteByte(0)
	}
	buffer.WriteByte(0)

	return buffer.String()
}

func (w *SelectInput) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountCombobox(parent.Handle, w.serializeItems(),
		w.Value, w.Unset, w.Disabled,
		w.OnChange != nil, w.OnFocus != nil, w.OnBlur != nil)

	retval := &selectinputElement{
		Control:  Control{control},
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}

func (w *selectinputElement) OnChange(value int) {
	if w.onChange != nil {
		w.onChange(value)
	}
}

func (w *selectinputElement) OnFocus() {
	w.onFocus()
}

func (w *selectinputElement) OnBlur() {
	w.onBlur()
}

func (w *selectinputElement) Props() base.Widget {
	value := gtk.ComboboxValue(w.handle)
	unset := value < 0
	if unset {
		value = 0
	}

	return &SelectInput{
		Items:    w.propsItems(),
		Value:    value,
		Unset:    unset,
		Disabled: !gtk.WidgetSensitive(w.handle),
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *selectinputElement) propsItems() []string {
	count := gtk.ComboboxItemCount(w.handle)

	items := []string{}
	for i := uint(0); i < count; i++ {
		items = append(items, gtk.ComboboxItem(w.handle, i))
	}

	return items
}

func (w *selectinputElement) TakeFocus() bool {
	control := Control{gtk.ComboboxChild(w.handle)}
	return control.TakeFocus()
}

func (w *selectinputElement) updateProps(data *SelectInput) error {
	w.onChange = nil // temporarily break OnChange to prevent event
	gtk.ComboboxUpdate(w.handle, data.serializeItems(), data.Value, data.Unset, data.Disabled,
		data.OnChange != nil, data.OnFocus != nil, data.OnBlur != nil)

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	return nil
}
