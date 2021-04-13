// +build cocoa darwin,!gtk

package dialog

import (
	"guitest/goey/internal/cocoa"
)

type dialogImpl struct {
	parent *cocoa.Window
}
