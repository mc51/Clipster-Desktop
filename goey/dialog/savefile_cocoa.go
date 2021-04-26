// +build cocoa darwin,!gtk

package dialog

import (
	"clipster/goey/internal/cocoa"
)

func (m *SaveFile) show() (string, error) {
	retval := cocoa.SavePanel(m.parent, m.filename)
	return retval, nil
}

// WithParent sets the parent of the dialog box.
func (m *SaveFile) WithParent(parent *cocoa.Window) *SaveFile {
	m.parent = parent
	return m
}
