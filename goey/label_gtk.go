// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/gtk"
)

type labelElement struct {
	Control
}

func (w *Label) mount(parent base.Control) (base.Element, error) {
	handle := gtk.MountLabel(parent.Handle, w.Text)

	retval := &labelElement{Control: Control{handle}}
	gtk.RegisterWidget(handle, retval)

	return retval, nil
}

func (w *labelElement) Props() base.Widget {
	return &Label{
		Text: gtk.LabelText(w.handle),
	}
}

func (w *labelElement) updateProps(data *Label) error {
	gtk.LabelUpdate(w.handle, data.Text)
	return nil
}
