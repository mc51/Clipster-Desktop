package dialog

import (
	"errors"
	"strings"
)

// SaveFile is a builder to construct an open file dialog to the user.
type SaveFile struct {
	Dialog
	title    string
	filename string
	filters  []filter
}

// NewSaveFile initializes a new open file dialog.
// Use of the method SaveFileDialog on an existing Window is preferred, as the
// message can be set as a child of the top-level window.
func NewSaveFile() *SaveFile {
	return &SaveFile{title: "goey"}
}

// Show completes building of the message, and shows the message to the user.
func (m *SaveFile) Show() (string, error) {
	if m.err != nil {
		return "", m.err
	}

	return m.show()
}

// AddFilter adds a filter to the list of filename patterns that the user can
// select.
func (m *SaveFile) AddFilter(name, pattern string) *SaveFile {
	m.filters = append(m.filters, filter{name, pattern})
	return m
}

// WithFilename sets the current or default filename for the dialog.
func (m *SaveFile) WithFilename(filename string) *SaveFile {
	m.filename = filename
	return m
}

// WithTitle adds a title to the dialog.
func (m *SaveFile) WithTitle(text string) *SaveFile {
	text = strings.TrimSpace(text)
	if text == "" {
		m.err = errors.New("Invalid argument, 'text' cannot be empty in call to WithTitle")
	} else {
		m.title = text
	}
	return m
}
