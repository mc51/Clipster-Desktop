// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package dialog

import (
	"guitest/goey/internal/gtk"
)

func (m *SaveFile) show() (string, error) {
	dlg := gtk.MountSaveDialog(m.parent, m.title, m.filename)
	activeDialogForTesting = dlg
	defer func() {
		activeDialogForTesting = 0
		gtk.WidgetClose(dlg)
	}()

	for _, v := range m.filters {
		gtk.DialogAddFilter(dlg, v.name, v.pattern)
	}

	rc := gtk.DialogRun(dlg)
	if rc != gtk.DialogResponseAccept() {
		return "", nil
	}
	return gtk.DialogGetFilename(dlg), nil
}

// WithParent sets the parent of the dialog box.
func (m *SaveFile) WithParent(parent uintptr) *SaveFile {
	m.parent = parent
	return m
}
