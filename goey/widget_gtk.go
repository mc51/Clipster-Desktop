// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"time"

	"clipster/goey/base"
	"clipster/goey/internal/gtk"
	"clipster/goey/loop"
)

// Control is an opaque type used as a platform-specific handle to a control
// created using the platform GUI.  As an example, this will refer to a HWND
// when targeting Windows, but a *GtkWidget when targeting GTK.
//
// Unless developping new widgets, users should not need to use this type.
//
// Any method's on this type will be platform specific.
type Control struct {
	handle uintptr
}

// Close removes the element from the GUI, and frees any associated resources.
func (w *Control) Close() {
	if w.handle != 0 {
		gtk.WidgetClose(w.handle)
		w.handle = 0
	}
}

// Handle returns the platform-native handle for the control.
func (w *Control) Handle() uintptr {
	return w.handle
}

func (w *Control) OnDestroy() {
	w.handle = 0
}

// TakeFocus is a wrapper around GrabFocus.
func (w *Control) TakeFocus() bool {
	// Check that the control can grab focus
	if !gtk.WidgetCanFocus(w.handle) {
		return false
	}

	gtk.WidgetGrabFocus(w.handle)
	// Note sure why the call to sleep is required, but there may be a debounce
	// provided by the system.  Without this call to sleep, the controls never
	// get the focus events.
	time.Sleep(250 * time.Millisecond)
	return gtk.WidgetIsFocus(w.handle)
}

// TypeKeys sends events to the control as if the string was typed by a user.
func (w *Control) TypeKeys(text string) chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		time.Sleep(500 * time.Millisecond)
		for _, r := range text {
			err := loop.Do(func() error {
				gtk.WidgetSendKey(w.handle, uint(r), false)
				return nil
			})
			if err != nil {
				errs <- err
			}
			time.Sleep(50 * time.Millisecond)

			err = loop.Do(func() error {
				gtk.WidgetSendKey(w.handle, uint(r), true)
				return nil
			})
			if err != nil {
				errs <- err
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	return errs
}

// Layout determines the best size for an element that satisfies the
// constraints.
func (w *Control) Layout(bc base.Constraints) base.Size {
	if !bc.HasBoundedWidth() && !bc.HasBoundedHeight() {
		// No need to worry about breaking the constraints.  We can take as
		// much space as desired.
		width, height := gtk.WidgetNaturalSize(w.handle)
		// Dimensions may need to be increased to meet minimums.
		return bc.Constrain(base.Size{base.FromPixelsX(width), base.FromPixelsY(height)})
	}
	if !bc.HasBoundedHeight() {
		// No need to worry about height.  Find the width that best meets the
		// widgets preferred width.
		width1 := gtk.WidgetNaturalWidth(w.handle)
		width := bc.ConstrainWidth(base.FromPixelsX(width1))
		// Get the best height for this width.
		height := gtk.WidgetNaturalHeightForWidth(w.handle, width.PixelsX())
		// Height may need to be increased to meet minimum.
		return base.Size{width, bc.ConstrainHeight(base.FromPixelsY(height))}
	}

	// Not clear the following is the best general approach given GTK layout
	// model.
	width, height := gtk.WidgetNaturalSize(w.handle)
	return bc.Constrain(base.Size{base.FromPixelsX(width), base.FromPixelsY(height)})
}

// MinIntrinsicHeight returns the minimum height that this element requires
// to be correctly displayed.
func (w *Control) MinIntrinsicHeight(width base.Length) base.Length {
	if width != base.Inf {
		height := gtk.WidgetMinHeightForWidth(w.handle, width.PixelsX())
		return base.FromPixelsY(height)
	}
	height := gtk.WidgetMinHeight(w.handle)
	return base.FromPixelsY(height)
}

// MinIntrinsicWidth returns the minimum width that this element requires
// to be correctly displayed.
func (w *Control) MinIntrinsicWidth(base.Length) base.Length {
	width := gtk.WidgetMinWidth(w.handle)
	return base.FromPixelsX(width)
}

// SetBounds updates the position of the widget.
func (w *Control) SetBounds(bounds base.Rectangle) {
	pixels := bounds.Pixels()
	if pixels.Dx() <= 0 || pixels.Dy() <= 0 {
		panic("internal error.  zero width or zero height bounds for control")
	}
	gtk.WidgetSetBounds(w.handle, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}
