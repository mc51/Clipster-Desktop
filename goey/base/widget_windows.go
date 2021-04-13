package base

import (
	"github.com/lxn/win"
)

const (
	// PLATFORM specifies the GUI toolkit being used.
	PLATFORM = "windows"
)

// Control is an opaque type used as a platform-specific handle to a control
// created using the platform GUI.  As an example, this will refer to a HWND
// when targeting Windows, but a *GtkContainer when targeting GTK.
//
// Unless developing new widgets, users should not need to use this type.
//
// Any methods on this type will be platform specific.
type Control struct {
	HWnd win.HWND
}

// NativeElement contains platform-specific methods that all widgets
// must support on WIN32
type NativeElement interface {
	// SetOrder is called to ensure that windows appears in the correct order,
	// which is important for tab navigation.  Elements should call SetOrder
	// on their children to create a depth-first traversal of all controls
	// (i.e. HWND).  Controls should use SetWindowPos to update their order.
	SetOrder(previous win.HWND) win.HWND
}
