// +build cocoa darwin,!gtk

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/cocoa"
)

type imgElement struct {
	control *cocoa.ImageView
	width   base.Length
	height  base.Length
}

func (w *Img) mount(parent base.Control) (base.Element, error) {
	// Convert the image to an NSImage
	control, err := cocoa.NewImageView(parent.Handle, w.Image)
	if err != nil {
		return nil, err
	}

	retval := &imgElement{
		control: control,
		width:   w.Width,
		height:  w.Height,
	}

	return retval, nil
}

func (w *imgElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *imgElement) Props() base.Widget {
	return &Img{
		Image:  w.control.Image(),
		Width:  w.width,
		Height: w.height,
	}
}

func (w *imgElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *imgElement) updateProps(data *Img) error {
	err := w.control.SetImage(data.Image)
	if err != nil {
		return err
	}
	w.width, w.height = data.Width, data.Height
	return nil
}
