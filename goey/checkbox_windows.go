package goey

import (
	"unsafe"

	"clipster/goey/base"
	"github.com/lxn/win"
)

func (w *Checkbox) mount(parent base.Control) (base.Element, error) {
	// Create the control.
	const STYLE = win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.BS_CHECKBOX | win.BS_TEXT | win.BS_NOTIFY
	hwnd, text, err := createControlWindow(0, &button.className[0], w.Text, STYLE, parent.HWnd)
	if err != nil {
		return nil, err
	}
	if w.Value {
		win.SendMessage(hwnd, win.BM_SETCHECK, win.BST_CHECKED, 0)
	}
	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &button.oldWindowProc, checkboxWindowProc)

	retval := &checkboxElement{
		Control:  Control{hwnd},
		text:     text,
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type checkboxElement struct {
	Control
	text     []uint16
	onChange func(value bool)
	onFocus  func()
	onBlur   func()
}

func (w *checkboxElement) Click() {
	win.SendMessage(w.hWnd, win.BM_CLICK, 0, 0)
}

func (w *checkboxElement) Props() base.Widget {
	return &Checkbox{
		Text:     w.Control.Text(),
		Value:    win.SendMessage(w.hWnd, win.BM_GETCHECK, 0, 0) == win.BST_CHECKED,
		Disabled: !win.IsWindowEnabled(w.hWnd),
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *checkboxElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *checkboxElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 17 * DIP
}

func (w *checkboxElement) MinIntrinsicWidth(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	width, _ := w.CalcRect(w.text)
	return base.FromPixelsX(int(width) + 17)
}

func (w *checkboxElement) updateProps(data *Checkbox) error {
	w.SetText(data.Text)
	w.SetDisabled(data.Disabled)
	if data.Value {
		win.SendMessage(w.hWnd, win.BM_SETCHECK, win.BST_CHECKED, 0)
	} else {
		win.SendMessage(w.hWnd, win.BM_SETCHECK, win.BST_UNCHECKED, 0)
	}

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}

func checkboxWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		checkboxGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := checkboxGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := checkboxGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		// WM_COMMAND is sent to the parent, which will only forward certain
		// message.  This code should only ever see EN_UPDATE, but we will
		// still check.
		switch notification := win.HIWORD(uint32(wParam)); notification {
		case win.BN_CLICKED:
			// Need to process the click to update the checkbox.
			check := uintptr(win.BST_CHECKED)
			if win.SendMessage(hwnd, win.BM_GETCHECK, 0, 0) == win.BST_CHECKED {
				check = win.BST_UNCHECKED
			}
			win.SendMessage(hwnd, win.BM_SETCHECK, check, 0)

			// Callback
			if w := checkboxGetPtr(hwnd); w.onChange != nil {
				w.onChange(check == win.BST_CHECKED)
			}
		}
		return 0
	}

	return win.CallWindowProc(button.oldWindowProc, hwnd, msg, wParam, lParam)
}

func checkboxGetPtr(hwnd win.HWND) *checkboxElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*checkboxElement)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
