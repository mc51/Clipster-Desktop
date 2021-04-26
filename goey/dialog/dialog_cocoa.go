// +build cocoa darwin,!gtk

package dialog

import (
	"clipster/goey/internal/cocoa"
)

type dialogImpl struct {
	parent *cocoa.Window
}
