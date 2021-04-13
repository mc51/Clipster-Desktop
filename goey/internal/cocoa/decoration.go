package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import (
	"image/color"
	"unsafe"
)

// Decoration is a wrapper for a GDecoration.
type Decoration struct {
	View
	private int
}

func toColor(clr color.Color) C.nscolor_t {
	r, g, b, a := clr.RGBA()
	return C.nscolor_t{
		r: C.uint8_t(r >> 8),
		g: C.uint8_t(g >> 8),
		b: C.uint8_t(b >> 8),
		a: C.uint8_t(a >> 8),
	}
}

func NewDecoration(window *View, fill color.Color, stroke color.Color, rx, ry int) *Decoration {
	handle := C.decorationNew(unsafe.Pointer(window),
		toColor(fill), toColor(stroke),
		C.nssize_t{C.int32_t(rx), C.int32_t(ry)})
	return (*Decoration)(handle)
}

func (w *Decoration) Close() {
	C.viewClose(unsafe.Pointer(w))
}

func (w *Decoration) BorderRadius() (int, int) {
	size := C.decorationBorderRadius(unsafe.Pointer(w))
	return int(size.width), int(size.height)
}

func (w *Decoration) FillColor() color.RGBA {
	clr := C.decorationFillColor(unsafe.Pointer(w))
	return color.RGBA{uint8(clr.r), uint8(clr.g), uint8(clr.b), uint8(clr.a)}
}

func (w *Decoration) StrokeColor() color.RGBA {
	clr := C.decorationStrokeColor(unsafe.Pointer(w))
	return color.RGBA{uint8(clr.r), uint8(clr.g), uint8(clr.b), uint8(clr.a)}
}

func (w *Decoration) SetBorderRadius(x, y int) {
	radius := C.nssize_t{
		width:  C.int32_t(x),
		height: C.int32_t(y),
	}
	C.decorationSetBorderRadius(unsafe.Pointer(w), radius)
}

func (w *Decoration) SetFillColor(fill color.Color) {
	C.decorationSetFillColor(unsafe.Pointer(w), toColor(fill))
}

func (w *Decoration) SetStrokeColor(fill color.Color) {
	C.decorationSetStrokeColor(unsafe.Pointer(w), toColor(fill))
}
