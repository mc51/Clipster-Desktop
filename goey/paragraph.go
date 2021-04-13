package goey

import (
	"guitest/goey/base"
)

var (
	paragraphKind     = base.NewKind("guitest/goey.Paragraph")
	paragraphMaxWidth base.Length
)

// TextAlignment identifies the different types of text alignment that are possible.
type TextAlignment uint8

// Allowed values for text alignment for text in paragraphs.
const (
	JustifyLeft   TextAlignment = iota // Text aligned to the left (ragged right)
	JustifyCenter                      // Text aligned to the center
	JustifyRight                       // Text aligned to the right (ragged left)
	JustifyFull                        // Text justified so that both left and right are flush
)

// P describes a widget that contains significant text, which can reflow if necessary.
//
// For a short run of text, the widget will try to match the size of the text.
// For longer runs of text, the widget will try to keep the width between 20em
// and 80em.
type P struct {
	Text  string        // Text is the content of the paragraph
	Align TextAlignment // Align is the text alignment for the paragraph
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*P) Kind() *base.Kind {
	return &paragraphKind
}

// Mount creates a paragraph in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *P) Mount(parent base.Control) (base.Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*paragraphElement) Kind() *base.Kind {
	return &paragraphKind
}

func (w *paragraphElement) Layout(bc base.Constraints) base.Size {
	if bc.HasBoundedWidth() {
		width := bc.ConstrainWidth(w.maxReflowWidth())
		height := w.MinIntrinsicHeight(width)
		return base.Size{width, bc.ConstrainHeight(height)}
	}

	if bc.HasBoundedHeight() {
		width := w.minReflowWidth()
		height := w.MinIntrinsicHeight(width)
		if height <= bc.Max.Height {
			return base.Size{width, height}
		}
		width = w.maxReflowWidth()
		height = w.MinIntrinsicHeight(width)
		return base.Size{width, bc.ConstrainHeight(height)}
	}

	width := bc.ConstrainWidth(w.minReflowWidth())
	height := w.MinIntrinsicHeight(width)
	return base.Size{width, bc.ConstrainHeight(height)}
}

func (w *paragraphElement) minReflowWidth() base.Length {
	if paragraphMaxWidth == 0 {
		w.measureReflowLimits()
	}
	// Get a minimum width of 20em compared to a max of 80em
	return paragraphMaxWidth / 4
}

func (w *paragraphElement) maxReflowWidth() base.Length {
	if paragraphMaxWidth == 0 {
		w.measureReflowLimits()
	}
	return paragraphMaxWidth
}

func (w *paragraphElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*P))
}
