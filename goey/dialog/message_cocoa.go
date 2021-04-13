// +build cocoa darwin,!gtk

package dialog

import (
	"guitest/goey/internal/cocoa"
)

func (m *Message) show() error {
	cocoa.MessageDialog(m.parent, m.text, m.title, byte(m.icon))
	return nil
}

func (m *Message) withError() {
	m.icon = 'e'
}

func (m *Message) withWarn() {
	m.icon = 'w'
}

func (m *Message) withInfo() {
	m.icon = 'i'
}

// WithParent sets the parent of the dialog box.
func (m *Message) WithParent(parent *cocoa.Window) *Message {
	m.parent = parent
	return m
}
