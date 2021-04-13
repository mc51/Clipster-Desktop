package dialog

import (
	"errors"
	"strings"
)

// OpenFile is a builder to construct an open file dialog to the user.
type OpenFile struct {
	Dialog
	title    string
	filename string
	filters  []filter
}

type filter struct {
	name    string
	pattern string
}

// NewOpenFile initializes a new open file dialog.
// Use of the method OpenFileDialog on an existing Window is preferred, as the
// message can be set as a child of the top-level window.
func NewOpenFile() *OpenFile {
	return &OpenFile{title: "goey"}
}

// Show completes building of the message, and shows the message to the user.
func (m *OpenFile) Show() (string, error) {
	if m.err != nil {
		return "", m.err
	}

	return m.show()
}

// AddFilter adds a filter to the list of filename patterns that the user can
// select.
func (m *OpenFile) AddFilter(name, pattern string) *OpenFile {
	m.filters = append(m.filters, filter{name, pattern})
	return m
}

// WithFilename sets the current or default filename for the dialog.
func (m *OpenFile) WithFilename(filename string) *OpenFile {
	m.filename = filename
	return m
}

// WithTitle adds a title to the dialog.
func (m *OpenFile) WithTitle(text string) *OpenFile {
	text = strings.TrimSpace(text)
	if text == "" {
		m.err = errors.New("Invalid argument, 'text' cannot be empty in call to WithTitle")
	} else {
		m.title = text
	}
	return m
}
