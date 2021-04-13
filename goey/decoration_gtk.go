// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"image/color"

	"guitest/goey/base"
	"guitest/goey/internal/gtk"
)

func toRGBA(c color.RGBA) uint {
	ret := uint32(c.R) | (uint32(c.G) << 8) | (uint32(c.B) << 16) | (uint32(c.A) << 24)
	return uint(ret)
}

func toColor(c uint) color.RGBA {
	return color.RGBA{
		R: uint8(c & 0xff),
		G: uint8((c >> 8) & 0xff),
		B: uint8((c >> 16) & 0xff),
		A: uint8((c >> 24) & 0xff),
	}
}

func (w *Decoration) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountDecoration(parent.Handle, toRGBA(w.Fill), toRGBA(w.Stroke), w.Radius.PixelsX())

	retval := &decorationElement{
		Control: Control{control},
		parent:  parent,
		radius:  w.Radius,
		insets:  w.Insets,
	}
	gtk.RegisterWidget(control, retval)

	child, err := base.Mount(parent, w.Child)
	if err != nil {
		gtk.WidgetClose(control)
		return nil, err
	}
	retval.child = child

	return retval, nil
}

type decorationElement struct {
	Control
	parent base.Control

	radius    base.Length
	insets    Insets
	child     base.Element
	childSize base.Size
}

func (w *decorationElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
	if w.handle != 0 {
		w.Control.Close()
	}
}

func (w *decorationElement) props() *Decoration {
	return &Decoration{
		Fill:   toColor(gtk.DecorationFill(w.handle)),
		Stroke: toColor(gtk.DecorationStroke(w.handle)),
		Insets: w.insets,
		Radius: base.FromPixelsX(gtk.DecorationRadius(w.handle)),
	}
}

func (w *decorationElement) SetBounds(bounds base.Rectangle) {
	w.Control.SetBounds(bounds)

	// The DPI informatio may not have been available (or up to date) when the
	// component was mounted.  Update the radius.
	gtk.DecorationSetRadius(w.handle, w.radius.PixelsX())

	bounds.Min.X += w.insets.Left
	bounds.Min.Y += w.insets.Top
	bounds.Max.X -= w.insets.Right
	bounds.Max.Y -= w.insets.Bottom
	w.child.SetBounds(bounds)
}

func (w *decorationElement) updateProps(data *Decoration) error {
	gtk.DecorationUpdate(w.handle, toRGBA(data.Fill), toRGBA(data.Stroke), data.Radius.PixelsX())
	w.radius = data.Radius

	child, err := base.DiffChild(w.parent, w.child, data.Child)
	if err != nil {
		return err
	}
	w.child = child
	w.insets = data.Insets

	return nil
}
