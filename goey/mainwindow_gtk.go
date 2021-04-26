// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"image"

	"clipster/goey/base"
	"clipster/goey/dialog"
	"clipster/goey/internal/gtk"
	"clipster/goey/loop"
)

type windowImpl struct {
	handle                  uintptr
	scroll                  uintptr
	layout                  uintptr
	child                   base.Element
	horizontalScroll        bool
	horizontalScrollVisible bool
	verticalScroll          bool
	verticalScrollVisible   bool
	onClosing               func() bool
	iconPix                 []byte
}

func newWindow(title string, child base.Widget) (*Window, error) {
	// Create a new GTK window
	window := gtk.MountWindow(title)
	loop.AddLockCount(1)

	retval := &Window{windowImpl{
		handle: window,
		scroll: gtk.WindowScrolledWindow(window),
		layout: gtk.WindowLayout(window),
	}}
	gtk.RegisterWidget(window, retval)
	gtk.WindowSetDefaultSize(func() (uintptr, int, int) {
		w, h := sizeDefaults()
		return window, int(w), int(h)
	}())

	return retval, nil
}

func (w *windowImpl) OnDestroy() {
	// Clear handle from the struct so that we dont' risk pointing to a
	// non existent window.
	w.handle = 0
	w.scroll = 0
	w.layout = 0
	// Release lock count on the GUI event loop.
	loop.AddLockCount(-1)

}

func (w *windowImpl) OnDeleteEvent() bool {
	if w.onClosing == nil {
		return false
	}
	return w.onClosing()
}

func (w *windowImpl) onSize() {
	w.OnSizeAllocate(gtk.WindowSize(w.handle))
}

func (w *windowImpl) OnSizeAllocate(width, height int) {
	if w.child == nil {
		return
	}

	// Update the global DPI
	base.DPI.X, base.DPI.Y = 96, 96

	clientSize := base.Size{base.FromPixelsX(width), base.FromPixelsY(height)}
	size := w.layoutChild(clientSize)
	if w.horizontalScroll && w.verticalScroll {
		// Show scroll bars if necessary.
		w.showScrollV(size.Height, clientSize.Height)
		ok := w.showScrollH(size.Width, clientSize.Width)
		// Adding horizontal scroll take vertical space, so we need to check
		// again for vertical scroll.
		if ok {
			_, height := gtk.WindowSize(w.handle)
			w.showScrollV(size.Height, base.FromPixelsY(height))
		}
	} else if w.verticalScroll {
		// Show scroll bars if necessary.
		ok := w.showScrollV(size.Height, clientSize.Height)
		if ok {
			width, height := gtk.WindowSize(w.handle)
			clientSize := base.Size{base.FromPixelsX(width), base.FromPixelsY(height)}
			size = w.layoutChild(clientSize)
		}
	} else if w.horizontalScroll {
		// Show scroll bars if necessary.
		ok := w.showScrollH(size.Width, clientSize.Width)
		if ok {
			width, height := gtk.WindowSize(w.handle)
			clientSize := base.Size{base.FromPixelsX(width), base.FromPixelsY(height)}
			size = w.layoutChild(clientSize)
		}
	}
	gtk.WindowSetLayoutSize(w.handle, uint(size.Width.PixelsX()), uint(size.Height.PixelsY()))
	bounds := base.Rectangle{
		base.Point{}, base.Point{size.Width, size.Height},
	}
	w.child.SetBounds(bounds)
}

func (w *windowImpl) control() base.Control {
	return base.Control{w.layout}
}

func (w *windowImpl) close() {
	if w.handle != 0 {
		gtk.WidgetClose(w.handle)
		w.handle = 0
	}
}

func (w *windowImpl) message(m *dialog.Message) {
	title := gtk.WindowTitle(w.handle)
	m.WithTitle(title)
	m.WithParent(w.handle)
}

func (w *windowImpl) openfiledialog(m *dialog.OpenFile) {
	m.WithParent(w.handle)
}

func (w *windowImpl) savefiledialog(m *dialog.SaveFile) {
	m.WithParent(w.handle)
}

// Screenshot returns an image of the window, as displayed on screen.
func (w *windowImpl) Screenshot() (image.Image, error) {
	pix, hasAlpha, width, height, stride := gtk.WindowScreenshot(w.handle)

	if hasAlpha {
		return &image.RGBA{
			Pix:    pix,
			Stride: stride,
			Rect:   image.Rect(0, 0, width, height),
		}, nil
	}

	newpix := make([]byte, height*width*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newpix[y*width*4+x*4+0] = pix[y*stride+x*3+0]
			newpix[y*width*4+x*4+1] = pix[y*stride+x*3+1]
			newpix[y*width*4+x*4+2] = pix[y*stride+x*3+2]
			newpix[y*width*4+x*4+3] = 0xff
		}
	}

	// Note:  stride of the new image data does not match data returned
	// from Pixbuf.
	return &image.RGBA{
		Pix:    newpix,
		Stride: width * 4,
		Rect:   image.Rect(0, 0, width, height),
	}, nil
}

func (w *windowImpl) setChildPost() {
	// Redo the layout so the children are placed.
	if w.child != nil {
		// Update the global DPI
		base.DPI.X, base.DPI.Y = 96, 96

		// Constrain window size
		w.updateWindowMinSize()
		// Properties may have changed sizes, so we need to do layout.
		w.onSize()
	} else {
		// Ensure that the scrollbars are hidden.
		gtk.WindowShowScrollbars(w.handle, false, false)
	}
}

func (w *windowImpl) setScroll(horz, vert bool) {
	w.horizontalScroll = horz
	w.verticalScroll = vert
	// Hide the scrollbars as a reset
	gtk.WindowShowScrollbars(w.handle, false, false)
	w.horizontalScrollVisible = false
	w.verticalScrollVisible = false
	// Redo layout to account for new box constraints, and show
	// scrollbars if necessary
	w.onSize()
}

func (w *windowImpl) show() {
	gtk.WindowShow(w.handle)
}

func (w *windowImpl) showScrollH(width base.Length, clientWidth base.Length) bool {
	if width > clientWidth {
		if !w.horizontalScrollVisible {
			// Show the scrollbar
			gtk.WindowShowScrollbars(w.handle, true, w.verticalScrollVisible)
			w.horizontalScrollVisible = true
			return true
		}
	} else if w.horizontalScrollVisible {
		// Remove the scroll bar
		gtk.WindowShowScrollbars(w.handle, false, w.verticalScrollVisible)
		w.horizontalScrollVisible = false
		return true
	}

	return false
}

func (w *windowImpl) showScrollV(height base.Length, clientHeight base.Length) bool {
	if height > clientHeight {
		if !w.verticalScrollVisible {
			// Show the scrollbar
			gtk.WindowShowScrollbars(w.handle, w.horizontalScrollVisible, true)
			w.verticalScrollVisible = true
			return true
		}
	} else if w.verticalScrollVisible {
		// Remove the scroll bar
		gtk.WindowShowScrollbars(w.handle, w.horizontalScrollVisible, false)
		w.verticalScrollVisible = false
		return true
	}

	return false
}

func (w *windowImpl) setIcon(img image.Image) error {
	if img == nil {
		gtk.WindowSetIcon(w.handle, nil, 0, 0, 0)
		w.iconPix = nil
	}

	rgba := imageToRGBA(img)
	gtk.WindowSetIcon(w.handle, &rgba.Pix[0], rgba.Rect.Dx(), rgba.Rect.Dy(), rgba.Stride)
	w.iconPix = rgba.Pix
	return nil
}

func (w *windowImpl) setOnClosing(callback func() bool) {
	w.onClosing = callback
}

func (w *windowImpl) setTitle(value string) error {
	gtk.WindowSetTitle(w.handle, value)
	return nil
}

func (w *windowImpl) title() (string, error) {
	return gtk.WindowTitle(w.handle), nil
}

func (w *windowImpl) updateWindowMinSize() {

	// Determine the extra width and height required for borders, title bar,
	// and scrollbars
	dx, dy := 0, 0
	if w.verticalScroll {
		dx += int(gtk.WindowVScrollbarWidth(w.handle))
	}
	if w.horizontalScroll {
		dy += int(gtk.WindowHScrollbarHeight(w.handle))
	}

	// If there is no child, then we just need enough space for the window chrome.
	if w.child == nil {
		gtk.WidgetSetSizeRequest(w.handle, dx, dy)
		return
	}

	request := image.Point{}
	// Determine the minimum size (in pixels) for the child of the window
	if w.horizontalScroll && w.verticalScroll {
		width := w.child.MinIntrinsicWidth(base.Inf)
		height := w.child.MinIntrinsicHeight(base.Inf)
		request.X = width.PixelsX() + dx
		request.Y = height.PixelsY() + dy
	} else if w.horizontalScroll {
		height := w.child.MinIntrinsicHeight(base.Inf)
		size := w.child.Layout(base.TightHeight(height))
		request.X = size.Width.PixelsX() + dx
		request.Y = height.PixelsY() + dy
	} else if w.verticalScroll {
		width := w.child.MinIntrinsicWidth(base.Inf)
		size := w.child.Layout(base.TightWidth(width))
		request.X = width.PixelsX() + dx
		request.Y = size.Height.PixelsY() + dy
	} else {
		width := w.child.MinIntrinsicWidth(base.Inf)
		height := w.child.MinIntrinsicHeight(base.Inf)
		size1 := w.child.Layout(base.TightWidth(width))
		size2 := w.child.Layout(base.TightHeight(height))
		request.X = max(width, size2.Width).PixelsX() + dx
		request.Y = max(height, size1.Height).PixelsY() + dy
	}

	// If scrolling is enabled for either direction, we can relax the
	// minimum window size.  These limits are fairly arbitrary, but we do need to
	// leave enough space for the scroll bars.
	if limit := (120 * DIP).PixelsX(); w.horizontalScroll && request.X > limit {
		request.X = limit
	}
	if limit := (120 * DIP).PixelsY(); w.verticalScroll && request.Y > limit {
		request.Y = limit
	}

	gtk.WidgetSetSizeRequest(w.handle, request.X, request.Y)
}
