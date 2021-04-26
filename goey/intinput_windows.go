package goey

import (
	"strconv"
	"syscall"
	"unsafe"

	"clipster/goey/base"
	win2 "clipster/goey/internal/syscall"
	"github.com/lxn/win"
)

var (
	intinput struct {
		className     []uint16
		oldWindowProc uintptr
	}
)

func init() {
	intinput.className = []uint16{'m', 's', 'c', 't', 'l', 's', '_', 'u', 'p', 'd', 'o', 'w', 'n', '3', '2', 0}
}

func (w *IntInput) mountUpDown(parent base.Control) (win.HWND, error) {
	// Range for the updown control is is only int32, not int64.
	if !w.useUpDownControl() {
		return 0, nil
	}

	hwnd, _, err := createControlWindow(win.WS_EX_LEFT|win.WS_EX_LTRREADING,
		&intinput.className[0],
		"",
		win.WS_CHILDWINDOW|win.WS_VISIBLE|win.UDS_SETBUDDYINT|win.UDS_ARROWKEYS|win.UDS_HOTTRACK|win.UDS_NOTHOUSANDS,
		parent.HWnd)
	if err != nil {
		return 0, err
	}
	win.SendMessage(hwnd, win.UDM_SETRANGE32, uintptr(w.Min), uintptr(w.Max))
	win.SendMessage(hwnd, win.UDM_SETPOS32, 0, uintptr(w.Value))

	return hwnd, nil
}

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_AUTOHSCROLL | win.ES_NUMBER)
	if w.OnEnterKey != nil {
		style = style | win.ES_MULTILINE
	}
	hwnd, _, err := createControlWindow(win.WS_EX_CLIENTEDGE, &edit.className[0], strconv.FormatInt(w.Value, 10), style, parent.HWnd)
	if err != nil {
		return nil, err
	}

	// Create the updown control.
	hwndUpDown, err := w.mountUpDown(parent)
	if err != nil {
		return nil, err
	}

	if w.Disabled {
		win.EnableWindow(hwnd, false)
		if hwndUpDown != 0 {
			win.EnableWindow(hwndUpDown, false)
		}
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
	retval := &intinputElement{
		Control:    Control{hwnd},
		hwndUpDown: hwndUpDown,
		min:        w.Min,
		max:        w.Max,
		onChange:   w.OnChange,
		onFocus:    w.OnFocus,
		onBlur:     w.OnBlur,
		onEnterKey: w.OnEnterKey,
	}

	// Link the control back to Go for event handling
	if hwndUpDown != 0 {
		win.SendMessage(hwndUpDown, win.UDM_SETBUDDY, uintptr(hwnd), 0)
		subclassWindowProcedure(hwnd, &intinput.oldWindowProc, intinputWindowProc)
	} else {
		subclassWindowProcedure(hwnd, &edit.oldWindowProc, intinputWindowProc)
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

func (w *IntInput) useUpDownControl() bool {
	// Range for the updown control is is only int32, not int64.
	// Need to make sure that we can properly set the range for the updown
	// control in order to include it in the GUI.
	return w.Min >= -2147483648 && w.Max <= 2147483647
}

type intinputElement struct {
	Control
	hwndUpDown win.HWND

	min        int64
	max        int64
	onChange   func(int64)
	onFocus    func()
	onBlur     func()
	onEnterKey func(int64)
}

func (w *intinputElement) Close() {
	if w.hwndUpDown != 0 {
		win.DestroyWindow(w.hwndUpDown)
		w.hwndUpDown = 0
	}
	if w.hWnd != 0 {
		win.DestroyWindow(w.hWnd)
		w.hWnd = 0
	}
}

func (w *intinputElement) getClampedValue() (int64, error) {
	// Get the text from the control, and convert text to an integer
	i, err := strconv.ParseInt(win2.GetWindowText(w.hWnd), 10, 64)
	if err != nil {
		return 0, err
	}
	// Clamp the value
	if i < w.min {
		i = w.min
		if w.hwndUpDown != 0 {
			win.SendMessage(w.hwndUpDown, win.UDM_SETPOS32, 0, uintptr(i))
			win.SendMessage(w.hWnd, win.EM_SETSEL, 0, 0x7fff)
		}
	} else if i > w.max {
		i = w.max
		if w.hwndUpDown != 0 {
			win.SendMessage(w.hwndUpDown, win.UDM_SETPOS32, 0, uintptr(i))
			win.SendMessage(w.hWnd, win.EM_SETSEL, 0, 0x7fff)
		}
	}

	return i, nil
}

func (w *intinputElement) thunkOnChange() {
	i, err := w.getClampedValue()
	if err != nil {
		// This case should not occur, as the control should prevent invalid
		// strings from being entered.
		// TODO:  What reporting should be done here?
		return
	}
	w.onChange(i)
}

func (w *intinputElement) thunkOnEnterKey() {
	i, err := w.getClampedValue()
	if err != nil {
		// This case should not occur, as the control should prevent invalid
		// strings from being entered.
		// TODO:  What reporting should be done here?
		return
	}
	w.onEnterKey(i)
}

func (w *intinputElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *intinputElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 * DIP
}

func (w *intinputElement) MinIntrinsicWidth(base.Length) base.Length {
	return 75 * DIP
}

func (w *intinputElement) Props() base.Widget {
	value := int64(0)
	if w.hwndUpDown != 0 {
		value = int64(win.SendMessage(w.hwndUpDown, win.UDM_GETPOS32, 0, 0))
	} else {
		value, _ = strconv.ParseInt(w.Control.Text(), 10, 64)
	}

	return &IntInput{
		Value:       value,
		Placeholder: propsPlaceholder(w.hWnd),
		Disabled:    !win.IsWindowEnabled(w.hWnd),
		Min:         w.min,
		Max:         w.max,
		OnChange:    w.onChange,
		OnFocus:     w.onFocus,
		OnBlur:      w.onBlur,
		OnEnterKey:  w.onEnterKey,
	}
}

func (w *intinputElement) SetBounds(bounds base.Rectangle) {
	buddyWidth := (23 * DIP) * 2 / 3

	if w.hwndUpDown == 0 {
		win.MoveWindow(w.hWnd, int32(bounds.Min.X.PixelsX()), int32(bounds.Min.Y.PixelsY()), int32(bounds.Dx().PixelsX()), int32(bounds.Dy().PixelsY()), false)
		return
	}

	if bounds.Dx() >= 4*buddyWidth {
		win.MoveWindow(w.hWnd, int32(bounds.Min.X.PixelsX()), int32(bounds.Min.Y.PixelsY()), int32((bounds.Dx() - buddyWidth).PixelsX()), int32(bounds.Dy().PixelsY()), false)
		win.MoveWindow(w.hwndUpDown, int32((bounds.Max.X - buddyWidth).PixelsX()), int32(bounds.Min.Y.PixelsY()), int32(buddyWidth.PixelsX()), int32(bounds.Dy().PixelsY()), false)
		win.ShowWindow(w.hwndUpDown, win.SW_SHOW)
	} else {
		win.MoveWindow(w.hWnd, int32(bounds.Min.X.PixelsX()), int32(bounds.Min.Y.PixelsY()), int32(bounds.Dx().PixelsX()), int32(bounds.Dy().PixelsY()), false)
		win.ShowWindow(w.hwndUpDown, win.SW_HIDE)
	}
}

func (w *intinputElement) SetOrder(previous win.HWND) win.HWND {
	if w.hwndUpDown != 0 {
		win.SetWindowPos(w.hwndUpDown, previous, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_NOREDRAW|0x400)
		previous = w.hwndUpDown
	}
	win.SetWindowPos(w.hWnd, previous, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_NOREDRAW|0x400)
	return w.hWnd
}

func (w *intinputElement) TakeFocus() bool {
	ok := w.Control.TakeFocus()
	if ok {
		win.SendMessage(w.hWnd, win.EM_SETSEL, 0, 0x7fff)
	}
	return ok
}

func (w *intinputElement) updateProps(data *IntInput) error {
	// Remove the updown control is the range is too large.
	if w.hwndUpDown != 0 && !data.useUpDownControl() {
		win.DestroyWindow(w.hwndUpDown)
		w.hwndUpDown = 0
	}

	text := strconv.FormatInt(data.Value, 10)
	if text != w.Text() {
		w.SetText(text)
	}
	if w.hwndUpDown != 0 {
		win.SendMessage(w.hwndUpDown, win.UDM_SETRANGE32, uintptr(data.Min), uintptr(data.Max))
		win.SendMessage(w.hwndUpDown, win.UDM_SETPOS32, 0, uintptr(data.Value))
	}
	err := updatePlaceholder(w.hWnd, data.Placeholder)
	if err != nil {
		return err
	}
	w.SetDisabled(data.Disabled)

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	w.onEnterKey = data.OnEnterKey

	return nil
}

func intinputWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		intinputGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := intinputGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := intinputGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_KEYDOWN:
		if wParam == win.VK_RETURN {
			if w := intinputGetPtr(hwnd); w.onEnterKey != nil {
				w.thunkOnEnterKey()
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
			if w := intinputGetPtr(hwnd); w.onChange != nil {
				w.thunkOnChange()
			}
		}
		return 0

	}

	if intinputGetPtr(hwnd).hwndUpDown != 0 {
		return win.CallWindowProc(intinput.oldWindowProc, hwnd, msg, wParam, lParam)
	}
	return win.CallWindowProc(edit.oldWindowProc, hwnd, msg, wParam, lParam)
}

func intinputGetPtr(hwnd win.HWND) *intinputElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*intinputElement)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
