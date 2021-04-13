package cocoa

/*
#include "cocoa.h"
*/
import "C"
import "unsafe"

type Control struct {
	View
	private int
}

func toBool(b bool) C.bool_t {
	if b {
		return 1
	}
	return 0
}

func (c *Control) IsEnabled() bool {
	rc := C.controlIsEnabled(unsafe.Pointer(c))
	return rc != 0
}

func (c *Control) IntrinsicContentSize() (int, int) {
	size := C.controlIntrinsicContentSize(unsafe.Pointer(c))
	return int(size.width), int(size.height)
}

func (c *Control) MakeFirstResponder() bool {
	return C.controlMakeFirstResponder(unsafe.Pointer(c)) != 0
}

func (c *Control) SetEnabled(enabled bool) {
	C.controlSetEnabled(unsafe.Pointer(c), toBool(enabled))
}

func (c *Control) SendKey(key uint) {
	C.controlSendKey(unsafe.Pointer(c), C.unsigned(key))
}
