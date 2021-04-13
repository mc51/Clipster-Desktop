// +build cocoa darwin,!gtk

package base

import (
	"guitest/goey/internal/cocoa"
)

const (
	// PLATFORM specifies the GUI toolkit being used.
	PLATFORM = "cocoa"
)

// Control is an opaque type used as a platform-specific handle to a control
// created using the platform GUI.  As an example, this will refer to a HWND
// when targeting Windows, but a *GtkContainer when targeting GTK.
//
// Unless developing new widgets, users should not need to use this type.
//
// Any methods on this type will be platform specific.
type Control struct {
	Handle *cocoa.View
}

// NativeElement contains platform-specific methods that all widgets
// must support on OSX.
type NativeElement interface {
}
