package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import (
	"path/filepath"
	"unsafe"
)

func MessageDialog(handle *Window, text string, title string, icon byte) {
	ctext := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctext))
	}()
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.messageDialog(unsafe.Pointer(handle), ctext, ctitle, C.char(icon))
}

func OpenPanel(handle *Window, filename string) string {
	var dir, base *C.char
	if filename != "" {
		dir = C.CString(filepath.Dir(filename))
		base = C.CString(filepath.Base(filename))
		defer func() {
			C.free(unsafe.Pointer(dir))
			C.free(unsafe.Pointer(base))
		}()
	}

	retval := C.openPanel(unsafe.Pointer(handle), dir, base)
	return C.GoString(retval)
}

func SavePanel(handle *Window, filename string) string {
	var dir, base *C.char
	if filename != "" {
		dir = C.CString(filepath.Dir(filename))
		base = C.CString(filepath.Base(filename))
		defer func() {
			C.free(unsafe.Pointer(dir))
			C.free(unsafe.Pointer(base))
		}()
	}

	retval := C.savePanel(unsafe.Pointer(handle), dir, base)
	return C.GoString(retval)
}

func DialogSendKey(key uint) {
	C.dialogSendKey(C.unsigned(key))
}
