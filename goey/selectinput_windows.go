package goey

import (
	"clipster/goey/base"
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	comboboxClassName     = []uint16{'C', 'O', 'M', 'B', 'O', 'B', 'O', 'X', 0}
	oldComboboxWindowProc uintptr
)

func (w *SelectInput) mount(parent base.Control) (base.Element, error) {
	const STYLE = win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.CBS_DROPDOWNLIST
	hwnd, _, err := createControlWindow(win.WS_EX_CLIENTEDGE, &comboboxClassName[0], "", STYLE, parent.HWnd)
	if err != nil {
		return nil, err
	}

	// Set the font for the window
	if hMessageFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, uintptr(hMessageFont), 0)
	}

	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Add items to the control
	longestString, err := selectinputAddItems(hwnd, w.Items)
	if err != nil {
		win.DestroyWindow(hwnd)
		return nil, err
	}
	if !w.Unset {
		win.SendMessage(hwnd, win.CB_SETCURSEL, uintptr(w.Value), 0)
	}

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &oldComboboxWindowProc, comboboxWindowProc)

	retval := &selectinputElement{
		Control:       Control{hwnd},
		onChange:      w.OnChange,
		onFocus:       w.OnFocus,
		onBlur:        w.OnBlur,
		longestString: longestString,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

func selectinputAddItems(hwnd win.HWND, items []string) (string, error) {
	longestString := ""
	for _, v := range items {
		text, err := syscall.UTF16PtrFromString(v)
		if err != nil {
			return "", err
		}
		win.SendMessage(hwnd, win.CB_ADDSTRING, 0, uintptr(unsafe.Pointer(text)))

		if len(v) > len(longestString) {
			longestString = v
		}
	}

	return longestString, nil
}

type selectinputElement struct {
	Control
	onChange func(value int)
	onFocus  func()
	onBlur   func()

	longestString  string
	preferredWidth base.Length
}

func (w *selectinputElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *selectinputElement) MinIntrinsicHeight(width base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 * DIP
}

func (w *selectinputElement) MinIntrinsicWidth(height base.Length) base.Length {
	if w.preferredWidth == 0 {
		text, err := syscall.UTF16FromString(w.longestString)
		if err != nil {
			w.preferredWidth = 75 * DIP
		} else {
			width, _ := w.CalcRect(text)
			w.preferredWidth = base.FromPixelsX(int(width)).Scale(13, 10)
		}
	}
	return w.preferredWidth
}

func (w *selectinputElement) Props() base.Widget {
	length := win.SendMessage(w.hWnd, win.CB_GETCOUNT, 0, 0)
	items := make([]string, int(length))
	for i := range items {
		buffer := [80]uint16{}
		length := win.SendMessage(w.hWnd, win.CB_GETLBTEXTLEN, uintptr(i), 0)
		if length > 79 {
			panic("not enough room")
		}
		win.SendMessage(w.hWnd, win.CB_GETLBTEXT, uintptr(i),
			uintptr(unsafe.Pointer(&buffer)))
		items[i] = syscall.UTF16ToString(buffer[:length])
	}
	value := win.SendMessage(w.hWnd, win.CB_GETCURSEL, 0, 0)
	unset := false
	// Depending on platform, the value may be either a 32-bit or a 64-bit
	// value, which somewhat complicates detecting CB_ERR.  The following
	// test works in both cases.
	if int32(value) == -1 /*win.CB_ERR, but bug with extension*/ {
		value, unset = 0, true
	}
	return &SelectInput{
		Items:    items,
		Value:    int(value),
		Unset:    unset,
		Disabled: !win.IsWindowEnabled(w.hWnd),
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *selectinputElement) updateProps(data *SelectInput) error {
	// This is a brute force approach.  The list of items is probably unchanged
	// most of the time.
	win.SendMessage(w.hWnd, win.CB_RESETCONTENT, 0, 0)
	longestString, err := selectinputAddItems(w.hWnd, data.Items)
	if err != nil {
		return err
	}
	if !data.Unset {
		win.SendMessage(w.hWnd, win.CB_SETCURSEL, uintptr(data.Value), 0)
	}

	w.SetDisabled(data.Disabled)
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	w.longestString = longestString
	// Clear cache
	w.preferredWidth = 0

	return nil
}

func comboboxWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		selectinputGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := selectinputGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := selectinputGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		// WM_COMMAND is sent to the parent, which will only forward certain
		// message.  This code should only ever see CBN_SELCHANGE, but we will
		// still check.
		switch notification := win.HIWORD(uint32(wParam)); notification {
		case win.CBN_SELCHANGE:
			if w := selectinputGetPtr(hwnd); w.onChange != nil {
				cursel := win.SendMessage(hwnd, win.CB_GETCURSEL, 0, 0)
				w.onChange(int(cursel))
			}
		}
		// defer to old window proc
	}

	return win.CallWindowProc(oldComboboxWindowProc, hwnd, msg, wParam, lParam)
}

func selectinputGetPtr(hwnd win.HWND) *selectinputElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*selectinputElement)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
