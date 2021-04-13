package dialog

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func (m *OpenFile) show() (string, error) {
	title, err := syscall.UTF16PtrFromString(m.title)
	if err != nil {
		return "", err
	}

	ofn := win.OPENFILENAME{
		LStructSize: uint32(unsafe.Sizeof(win.OPENFILENAME{})),
		HwndOwner:   m.hWnd,
		LpstrFilter: buildFilterString(m.filters),
		LpstrTitle:  title,
		Flags:       win.OFN_PATHMUSTEXIST | win.OFN_FILEMUSTEXIST,
	}

	filename := [1024]uint16{}
	buffer, err := buildFileString(&ofn, filename[:], m.filename)
	if err != nil {
		return "", err
	}

	rc := win.GetOpenFileName(&ofn)
	if !rc {
		if err := win.CommDlgExtendedError(); err != 0 {
			return "", fmt.Errorf("call to GetOpenFileName failed with code %x", err)
		}
		return "", nil
	}
	return syscall.UTF16ToString(buffer), nil
}

func buildFilterString(filters []filter) *uint16 {
	// If there are no filters, we want to return a nil pointer.
	// This will let windows select appropriate default behavior.
	if len(filters) == 0 {
		return nil
	}

	// See documentation for OPENFILENAME structures, but we build a single
	// buffer with pairs of null-terminated strings.
	buffer := make([]uint16, 0, 1024)
	for _, v := range filters {
		tmp, _ := syscall.UTF16FromString(v.name)
		buffer = append(buffer, tmp...)
		tmp, _ = syscall.UTF16FromString(v.pattern)
		buffer = append(buffer, tmp...)
	}

	// Final double null-terminated string marks end of buffer.
	buffer = append(buffer, 0, 0)
	return &buffer[0]
}

func buildFileString(ofn *win.OPENFILENAME, buffer []uint16, filename string) ([]uint16, error) {
	if filename == "" {
		ofn.LpstrFile = &buffer[0]
		ofn.NMaxFile = uint32(cap(buffer))
		return buffer, nil
	}

	tmp, err := syscall.UTF16FromString(filepath.FromSlash(filename))
	if err != nil {
		return nil, err
	}

	// Copy the filename (with path) into the buffer, extending its size if
	// necessary.
	buffer = append(buffer[:0], tmp...)

	// In case the filename is now longer
	ofn.LpstrFile = &buffer[0]
	ofn.NMaxFile = uint32(cap(buffer))
	return buffer[:cap(buffer)], nil
}

// WithOwner sets the owner of the dialog box.
func (m *OpenFile) WithOwner(hwnd win.HWND) *OpenFile {
	m.hWnd = hwnd
	return m
}
