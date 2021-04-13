package cocoa

/*
#include "cocoa.h"
*/
import "C"
import "unsafe"

// View is a wrapper for a NSView.
type View struct {
	private int
}

// Close is a a wrapper for removeFromSuperview + release.
func (c *View) Close() {
	C.viewClose(unsafe.Pointer(c))
}

// SetFrame is a wrapper for the setFrame message.
func (c *View) SetFrame(x, y, dx, dy int) {
	C.viewSetFrame(unsafe.Pointer(c), C.int(x), C.int(y), C.int(dx), C.int(dy))
}
