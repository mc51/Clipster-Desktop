package goey

import (
	"image"

	"guitest/goey/base"
)

var (
	imgKind = base.NewKind("guitest/goey.Img")
)

// Img describes a widget that contains a bitmap image.
//
// The size of the control depends on the value of Width and Height.
// The fields Width and Height may be left uninitialized, in which case they
// will be modified in-place.  If both of these fields are left
// as zero, then the size will be calculated from the image's size assuming
// that its resolution is 92 DPI.  If only one dimension is zero, then it will
// be calculate to maintain the aspect ratio of the image.
type Img struct {
	Image         image.Image // Image to be displayed.
	Width, Height base.Length // Dimensions for the image (see notes on sizing).
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Img) Kind() *base.Kind {
	return &imgKind
}

// Mount creates an image control in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Img) Mount(parent base.Control) (base.Element, error) {
	// Fill in the height and width if they are left at zero.
	w.UpdateDimensions()

	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*imgElement) Kind() *base.Kind {
	return &imgKind
}

func (w *imgElement) Layout(bc base.Constraints) base.Size {
	// Determine ideal width.
	return bc.ConstrainAndAttemptToPreserveAspectRatio(base.Size{w.width, w.height})
}

func (w *imgElement) MinIntrinsicHeight(width base.Length) base.Length {
	if width < w.width {
		return w.height * width / w.width
	}
	return w.height
}

func (w *imgElement) MinIntrinsicWidth(height base.Length) base.Length {
	if height < w.height {
		return w.width * height / w.height
	}
	return w.width
}

// UpdateDimensions calculates default values for Width and Height if either
// or zero based on the image dimensions.  The member Image cannot be nil.
func (w *Img) UpdateDimensions() {
	if w.Width == 0 && w.Height == 0 {
		bounds := w.Image.Bounds()
		// Assume that images are at 92 pixels per inch
		w.Width = (1 * Inch).Scale(bounds.Dx(), 92)
		w.Height = (1 * Inch).Scale(bounds.Dy(), 92)
	} else if w.Width == 0 {
		bounds := w.Image.Bounds()
		w.Width = w.Height.Scale(bounds.Dx(), bounds.Dy())
	} else if w.Height == 0 {
		bounds := w.Image.Bounds()
		w.Height = w.Width.Scale(bounds.Dy(), bounds.Dx())
	}
}

func (w *imgElement) UpdateProps(data base.Widget) error {
	img := data.(*Img)

	// Fill in the height and width if they are left at zero.
	img.UpdateDimensions()
	// Forward to the platform-dependant code
	return w.updateProps(data.(*Img))
}
