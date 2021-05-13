// +build gtk linux darwin freebsd openbsd

package dialog

import (
	"clipster/goey/internal/gtk"
)

func (m *Message) show() error {
	dlg := gtk.MountMessageDialog(m.parent, m.title, m.icon, m.text)
	activeDialogForTesting = dlg
	defer func() {
		activeDialogForTesting = 0
		gtk.WidgetClose(dlg)
	}()

	gtk.DialogRun(dlg)
	return nil
}

func (m *Message) withError() {
	m.icon = gtk.MessageDialogWithError()
}

func (m *Message) withWarn() {
	m.icon = gtk.MessageDialogWithWarn()
}

func (m *Message) withInfo() {
	m.icon = gtk.MessageDialogWithInfo()
}

// WithParent sets the parent of the dialog box.
func (m *Message) WithParent(parent uintptr) *Message {
	m.parent = parent
	return m
}
