// +build gtk linux darwin freebsd openbsd

package dialog

import (
	"clipster/goey/internal/gtk"
)

func (m *OpenFile) show() (string, error) {
	dlg := gtk.MountOpenDialog(m.parent, m.title, m.filename)
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
func (m *OpenFile) WithParent(parent uintptr) *OpenFile {
	m.parent = parent
	return m
}
