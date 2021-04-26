package goey

import (
	"time"
	"unsafe"

	"clipster/goey/base"
	win2 "clipster/goey/internal/syscall"
	"github.com/lxn/win"
)

var (
	datetimepickClassName     []uint16
	oldDateTimePickWindowProc uintptr
)

func init() {
	datetimepickClassName = []uint16{'S', 'y', 's', 'D', 'a', 't', 'e', 'T', 'i', 'm', 'e', 'P', 'i', 'c', 'k', '3', '2', 0}
}

func (w *DateInput) systemTime() win.SYSTEMTIME {
	return win.SYSTEMTIME{
		WYear:   uint16(w.Value.Year()),
		WMonth:  uint16(w.Value.Month()),
		WDay:    uint16(w.Value.Day()),
		WHour:   uint16(w.Value.Hour()),
		WMinute: uint16(w.Value.Minute()),
		WSecond: uint16(w.Value.Second()),
	}
}

func (w *DateInput) mount(parent base.Control) (base.Element, error) {
	const STYLE = win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP
	hwnd, _, err := createControlWindow(0, &datetimepickClassName[0], "", STYLE, parent.HWnd)
	if err != nil {
		return nil, err
	}

	// Set the properties for the control
	st := w.systemTime()
	win.SendMessage(hwnd, win.DTM_SETSYSTEMTIME, win.GDT_VALID, uintptr(unsafe.Pointer(&st)))
	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &oldDateTimePickWindowProc, dateinputWindowProc)

	retval := &dateinputElement{
		Control:  Control{hwnd},
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type dateinputElement struct {
	Control
	onChange func(value time.Time)
	onFocus  func()
	onBlur   func()
}

func (w *dateinputElement) Layout(bc base.Constraints) base.Size {
	height := w.MinIntrinsicHeight(0)
	width := w.MinIntrinsicWidth(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *dateinputElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 * DIP
}

func (w *dateinputElement) MinIntrinsicWidth(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 75 * DIP
}

func (w *dateinputElement) Props() base.Widget {
	st := win.SYSTEMTIME{}
	win.SendMessage(w.hWnd, win.DTM_GETSYSTEMTIME, 0, uintptr(unsafe.Pointer(&st)))

	return &DateInput{
		Value: time.Date(int(st.WYear), time.Month(st.WMonth), int(st.WDay),
			int(st.WHour), int(st.WMinute), int(st.WSecond), 0, time.Local),
		Disabled: !win.IsWindowEnabled(w.hWnd),
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *dateinputElement) updateProps(data *DateInput) error {
	st := data.systemTime()
	win.SendMessage(w.hWnd, win.DTM_SETSYSTEMTIME, win.GDT_VALID, uintptr(unsafe.Pointer(&st)))

	w.SetDisabled(data.Disabled)
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	return nil
}

func dateinputWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		dateinputGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := dateinputGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := dateinputGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_NOTIFY:
		switch code := (*win.NMHDR)(unsafe.Pointer(lParam)).Code; code {
		case win.DTN_DATETIMECHANGE:
			if w := dateinputGetPtr(hwnd); w.onChange != nil {
				nmhdr := (*win.NMDATETIMECHANGE)(unsafe.Pointer(lParam))
				st := time.Date(int(nmhdr.St.WYear), time.Month(nmhdr.St.WMonth), int(nmhdr.St.WDay), int(nmhdr.St.WHour), int(nmhdr.St.WMinute), int(nmhdr.St.WSecond), 0, time.Local)
				w.onChange(st)
			}

		case win2.MCN_SELECT:
			nmhdr := (*win2.NMSELCHANGE)(unsafe.Pointer(lParam))
			win.SendMessage(hwnd, win.DTM_SETSYSTEMTIME, win.GDT_VALID, uintptr(unsafe.Pointer(&nmhdr.StSelStart)))
			win.SendMessage(hwnd, win2.DTM_CLOSEMONTHCAL, 0, 0)
		}
		return 0

	}

	return win.CallWindowProc(oldDateTimePickWindowProc, hwnd, msg, wParam, lParam)
}

func dateinputGetPtr(hwnd win.HWND) *dateinputElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*dateinputElement)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
