package goey

import (
	"syscall"
	"unsafe"

	"clipster/goey/base"
	win2 "clipster/goey/internal/syscall"
	"github.com/lxn/win"
)

var (
	edit struct {
		className     []uint16
		oldWindowProc uintptr
		emptyString   uint16
	}
)

func init() {
	edit.className = []uint16{'E', 'D', 'I', 'T', 0}
}

func (w *TextInput) mount(parent base.Control) (base.Element, error) {
	// Create the control.
	hwnd, _, err := createControlWindow(win.WS_EX_CLIENTEDGE, &edit.className[0], w.Value, w.style(), parent.HWnd)
	if err != nil {
		return nil, err
	}
	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Create placeholder, if required.
	if w.Placeholder != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(w.Placeholder)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		win.SendMessage(hwnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	}

	// Create the return value.
	retval := &textinputElement{textinputElementBase{
		Control:    Control{hwnd},
		onChange:   w.OnChange,
		onFocus:    w.OnFocus,
		onBlur:     w.OnBlur,
		onEnterKey: w.OnEnterKey,
	}}

	// Link the control back to Go for event handling
	subclassWindowProcedure(hwnd, &edit.oldWindowProc, textinputWindowProc)
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

func (w *TextInput) style() uint32 {
	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_AUTOHSCROLL)
	if w.Password {
		style = style | win.ES_PASSWORD
	}
	if w.ReadOnly {
		style = style | win.ES_READONLY
	}
	if w.OnEnterKey != nil {
		style = style | win.ES_MULTILINE
	}
	return style
}

type textinputElementBase struct {
	Control
	onChange   func(value string)
	onFocus    func()
	onBlur     func()
	onEnterKey func(value string)
}

type textinputElement struct {
	textinputElementBase
}

func (w *textinputElement) Props() base.Widget {
	return &TextInput{
		Value:       w.Control.Text(),
		Placeholder: propsPlaceholder(w.hWnd),
		Disabled:    !win.IsWindowEnabled(w.hWnd),
		Password:    win.SendMessage(w.hWnd, win.EM_GETPASSWORDCHAR, 0, 0) != 0,
		ReadOnly:    (win.GetWindowLong(w.hWnd, win.GWL_STYLE) & win.ES_READONLY) != 0,
		OnChange:    w.onChange,
		OnFocus:     w.onFocus,
		OnBlur:      w.onBlur,
		OnEnterKey:  w.onEnterKey,
	}
}

func propsPlaceholder(hWnd win.HWND) string {
	var buffer [80]uint16
	win.SendMessage(hWnd, win.EM_GETCUEBANNER, uintptr(unsafe.Pointer(&buffer[0])), 80)
	ndx := 0
	for i, v := range buffer {
		if v == 0 {
			ndx = i
			break
		}
	}
	return syscall.UTF16ToString(buffer[:ndx])
}

func (w *textinputElementBase) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *textinputElementBase) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 * DIP
}

func (w *textinputElementBase) MinIntrinsicWidth(base.Length) base.Length {
	// TODO
	return 75 * DIP
}

func (w *textinputElementBase) TakeFocus() bool {
	ok := w.Control.TakeFocus()
	if ok {
		win.SendMessage(w.hWnd, win.EM_SETSEL, 0, 0x7fff)
	}
	return ok
}

func updatePlaceholder(hWnd win.HWND, text string) error {
	// Update the control
	if text != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(text)
		if err != nil {
			return err
		}

		win.SendMessage(hWnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	} else {
		win.SendMessage(hWnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(&edit.emptyString)))
	}

	return nil
}

func (w *textinputElementBase) updateProps(data *TextInput) error {
	if data.Value != w.Text() {
		w.SetText(data.Value)
	}
	err := updatePlaceholder(w.hWnd, data.Placeholder)
	if err != nil {
		return err
	}
	w.SetDisabled(data.Disabled)
	if data.Password {
		// TODO:  ???
	} else {
		win.SendMessage(w.hWnd, win.EM_SETPASSWORDCHAR, 0, 0)
	}
	win.SendMessage(w.hWnd, win.EM_SETREADONLY, uintptr(win.BoolToBOOL(data.ReadOnly)), 0)

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	w.onEnterKey = data.OnEnterKey

	return nil
}

func textinputWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		textinputGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := textinputGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := textinputGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_KEYDOWN:
		if wParam == win.VK_RETURN {
			if w := textinputGetPtr(hwnd); w.onEnterKey != nil {
				w.onEnterKey(win2.GetWindowText(hwnd))
				return 0
			}
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		// WM_COMMAND is sent to the parent, which will only forward certain
		// message.  This code should only ever see EN_UPDATE, but we will
		// still check.
		switch notification := win.HIWORD(uint32(wParam)); notification {
		case win.EN_UPDATE:
			if w := textinputGetPtr(hwnd); w.onChange != nil {
				w.onChange(win2.GetWindowText(hwnd))
			}
		}
		return 0

	}

	return win.CallWindowProc(edit.oldWindowProc, hwnd, msg, wParam, lParam)
}

func textinputGetPtr(hwnd win.HWND) *textinputElementBase {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*textinputElementBase)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
