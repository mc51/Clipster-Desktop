package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Text is a wrapper for a NSText.
type Text struct {
	View
	private int
}

func NewText(window *View, title string) *Text {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.textNew(unsafe.Pointer(window), ctitle)
	return (*Text)(handle)
}

func (w *Text) Alignment() int {
	rc := C.textAlignment(unsafe.Pointer(w))
	return int(rc)
}

func (w *Text) EightyEms() int {
	csize := C.textEightyEms(unsafe.Pointer(w))
	return int(csize)
}

func (w *Text) MinHeight(width int) int {
	csize := C.textMinHeight(unsafe.Pointer(w), C.int(width))
	return int(csize)
}

func (w *Text) MinWidth() int {
	csize := C.textMinWidth(unsafe.Pointer(w))
	return int(csize)
}

func (w *Text) SetAlignment(align int) {
	C.textSetAlignment(unsafe.Pointer(w), C.int(align))
}

func (w *Text) SetText(text string) {
	ctext := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctext))
	}()

	C.textSetText(unsafe.Pointer(w), ctext)
}

func (w *Text) Text() string {
	ctext := C.textText(unsafe.Pointer(w))
	return C.GoString(ctext)
}
