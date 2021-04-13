// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"time"

	"guitest/goey/base"
	"guitest/goey/internal/gtk"
)

type dateinputElement struct {
	Control

	onChange func(time.Time)
	onFocus  func()
	onBlur   func()
}

func (w *DateInput) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountDateInput(parent.Handle,
		w.Value.Year(), uint(w.Value.Month()), uint(w.Value.Day()), w.Disabled,
		w.OnChange != nil, w.OnFocus != nil, w.OnBlur != nil)

	// Create the element
	retval := &dateinputElement{
		Control:  Control{control},
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}

func (w *dateinputElement) OnChange(value time.Time) {
	if w.onChange != nil {
		w.onChange(value)
	}
}

func (w *dateinputElement) OnFocus() {
	w.onFocus()
}

func (w *dateinputElement) OnBlur() {
	w.onBlur()
}

func (w *dateinputElement) Props() base.Widget {
	year := gtk.DateInputYear(w.handle)
	month := gtk.DateInputMonth(w.handle)
	day := gtk.DateInputDay(w.handle)

	return &DateInput{
		Value:    time.Date(year, time.Month(month), int(day), 0, 0, 0, 0, time.Local),
		Disabled: !gtk.WidgetSensitive(w.handle),
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}

}

func (w *dateinputElement) updateProps(data *DateInput) error {
	w.onChange = nil // temporarily break OnChange to prevent event
	gtk.DateInputUpdate(w.handle,
		data.Value.Year(), uint(data.Value.Month()), uint(data.Value.Day()), data.Disabled,
		data.OnChange != nil, data.OnFocus != nil, data.OnBlur != nil)

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}
