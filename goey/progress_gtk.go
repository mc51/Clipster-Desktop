// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/gtk"
)

type progressElement struct {
	Control
	min, max int
}

func (w *Progress) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountProgressbar(parent.Handle, float64(w.Value-w.Min)/float64(w.Max-w.Min))

	retval := &progressElement{
		Control: Control{control},
		min:     w.Min,
		max:     w.Max,
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}

func (w *progressElement) Props() base.Widget {
	if w.min == w.max {
		return &Progress{
			Value: w.min,
			Min:   w.min,
			Max:   w.max,
		}
	}

	return &Progress{
		Value: w.min + int(float64(w.max-w.min)*gtk.ProgressbarValue(w.handle)),
		Min:   w.min,
		Max:   w.max,
	}
}

func (w *progressElement) updateProps(data *Progress) error {
	w.min = data.Min
	w.max = data.Max
	gtk.ProgressbarUpdate(w.handle, float64(data.Value-data.Min)/float64(data.Max-data.Min))
	return nil
}
