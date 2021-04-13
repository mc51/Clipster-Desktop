// +build cocoa darwin,!gtk

package dialog

import (
	"guitest/goey/internal/cocoa"
)

func (m *OpenFile) show() (string, error) {
	retval := cocoa.OpenPanel(m.parent, m.filename)
	return retval, nil
}

// WithParent sets the parent of the dialog box.
func (m *OpenFile) WithParent(parent *cocoa.Window) *OpenFile {
	m.parent = parent
	return m
}
