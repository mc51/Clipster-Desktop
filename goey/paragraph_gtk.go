// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/gtk"
)

type paragraphElement struct {
	Control
}

func (w *P) mount(parent base.Control) (base.Element, error) {
	handle := gtk.MountParagraph(parent.Handle, w.Text, byte(w.Align))

	retval := &paragraphElement{Control{handle}}
	gtk.RegisterWidget(handle, retval)

	return retval, nil
}

func (w *paragraphElement) Props() base.Widget {
	return &P{
		Text:  gtk.ParagraphText(w.handle),
		Align: TextAlignment(gtk.ParagraphAlign(w.handle)),
	}
}

func (w *paragraphElement) measureReflowLimits() {
	oldText := gtk.ParagraphText(w.handle)

	gtk.ParagraphSetText(w.handle, "mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm")
	width := gtk.WidgetMinWidth(w.handle)
	gtk.ParagraphSetText(w.handle, oldText)

	paragraphMaxWidth = base.FromPixelsX(width)
}

func (w *paragraphElement) MinIntrinsicHeight(width base.Length) base.Length {
	if width == base.Inf {
		width = w.maxReflowWidth()
	}

	height := gtk.WidgetMinHeightForWidth(w.handle, width.PixelsX())
	return base.FromPixelsY(height)
}

func (w *paragraphElement) MinIntrinsicWidth(height base.Length) base.Length {
	if height != base.Inf {
		panic("not implemented")
	}

	width := gtk.WidgetNaturalWidth(w.handle)
	return min(base.FromPixelsX(width), w.minReflowWidth())
}

func (w *paragraphElement) updateProps(data *P) error {
	gtk.ParagraphUpdate(w.handle, data.Text, byte(data.Align))
	return nil
}
