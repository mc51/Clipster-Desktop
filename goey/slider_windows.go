package goey

import (
	"unsafe"

	"guitest/goey/base"
	"github.com/lxn/win"
)

var (
	slider struct {
		className     []uint16
		oldWindowProc uintptr
	}
)

func init() {
	slider.className = []uint16{'m', 's', 'c', 't', 'l', 's', '_', 't', 'r', 'a', 'c', 'k', 'b', 'a', 'r', '3', '2', 0}
}

func (w *Slider) mount(parent base.Control) (base.Element, error) {
	// This value should be as large as possible to maximize the resolution
	// of the slider.  However, if it is too large, then it will trip a bug
	// on windows, causing very high CPU usage.
	const RANGEMAX = 0xffffff

	const TBS_HORZ = 0x0000
	const TBS_AUTOTICKS = 0x0001
	const TBM_SETTICFREQ = win.WM_USER + 20
	const TBM_SETPAGESIZE = win.WM_USER + 21
	const TBM_SETLINESIZE = win.WM_USER + 23

	// Create the control
	const STYLE = win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | TBS_HORZ | TBS_AUTOTICKS
	hwnd, _, err := createControlWindow(0, &slider.className[0], "", STYLE, parent.HWnd)
	if err != nil {
		return nil, err
	}
	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}
	win.SendMessage(hwnd, win.TBM_SETRANGEMAX, win.FALSE, RANGEMAX)
	win.SendMessage(hwnd, TBM_SETLINESIZE, win.FALSE, RANGEMAX/100)
	win.SendMessage(hwnd, TBM_SETPAGESIZE, win.FALSE, RANGEMAX/16)
	win.SendMessage(hwnd, TBM_SETTICFREQ, win.FALSE, RANGEMAX/8)
	currentValue := sliderToQuantized(w.Value, w.Min, w.Max)
	win.SendMessage(hwnd, win.TBM_SETPOS, win.TRUE, currentValue)

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &slider.oldWindowProc, sliderWindowProc)

	retval := &sliderElement{
		Control:      Control{hwnd},
		currentValue: currentValue,
		min:          w.Min,
		max:          w.Max,
		onChange:     w.OnChange,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type sliderElement struct {
	Control
	currentValue uintptr
	min, max     float64

	onChange func(float64)
	onFocus  func()
	onBlur   func()
}

func (w *sliderElement) Layout(bc base.Constraints) base.Size {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	width := w.MinIntrinsicWidth(0)
	if bc.Max.Width > 355*DIP {
		width = 355 * DIP
	}
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *sliderElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 24 * DIP
}

func (w *sliderElement) MinIntrinsicWidth(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 160 * DIP
}

func (w *sliderElement) Props() base.Widget {
	currentValue := win.SendMessage(w.hWnd, win.TBM_GETPOS, 0, 0)
	// We will get errors in testing because of rounding errors in conversion
	// of float to int and back.

	return &Slider{
		Value:    sliderFromQuantized(currentValue, w.min, w.max),
		Disabled: !win.IsWindowEnabled(w.hWnd),
		Min:      w.min,
		Max:      w.max,
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *sliderElement) updateProps(data *Slider) error {
	w.min = data.Min
	w.max = data.Max
	if newValue := sliderToQuantized(data.Value, w.min, w.max); newValue != w.currentValue {
		w.currentValue = newValue
		win.SendMessage(w.hWnd, win.TBM_SETPOS, win.TRUE, newValue)
	}
	win.EnableWindow(w.hWnd, !data.Disabled)
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	return nil
}

func sliderFromQuantized(value uintptr, min, max float64) float64 {
	// Perform the conversion.
	retval := min + float64(value)*(max-min)/0xffffff

	// Apply correction if necessary to account for precision of slider.
	// The implementation on windows has a limited resolution compared to
	// float64.
	if tmp := float64(int64(retval*8)) / 8; tmp != retval {
		retval = tmp
	}
	return retval
}

func sliderToQuantized(value, min, max float64) uintptr {
	return uintptr((value-min)/(max-min)*0xffffff + 0.5)
}

func sliderWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		sliderGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := sliderGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := sliderGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_HSCROLL:
		// When the event is to end the scroll, the new position is not sent.
		// Skip these messages.
		if w := sliderGetPtr(hwnd); w.onChange != nil {
			if code := wParam & 0xffff; code >= win.SB_LINELEFT && code <= win.SB_PAGERIGHT {
				ret := win.CallWindowProc(slider.oldWindowProc, hwnd, msg, wParam, lParam)
				w.currentValue = win.SendMessage(hwnd, win.TBM_GETPOS, 0, 0)
				w.onChange(sliderFromQuantized(w.currentValue, w.min, w.max))
				return ret
			} else if code == win.SB_LEFT {
				win.SendMessage(hwnd, win.TBM_SETPOS, win.TRUE, 0)
				w.currentValue = 0
				w.onChange(w.min)
			} else if code == win.SB_RIGHT {
				win.SendMessage(hwnd, win.TBM_SETPOS, win.TRUE, 0xffff)
				w.currentValue = 0xffff
				w.onChange(w.max)
			} else if code == win.SB_THUMBPOSITION {
				w.currentValue = wParam >> 16
				w.onChange(sliderFromQuantized(w.currentValue, w.min, w.max))
			} else if code == win.SB_THUMBTRACK {
				w.currentValue = wParam >> 16
				w.onChange(sliderFromQuantized(w.currentValue, w.min, w.max))
			}
		}
		// Defer to the old window proc
	}

	return win.CallWindowProc(slider.oldWindowProc, hwnd, msg, wParam, lParam)
}

func sliderGetPtr(hwnd win.HWND) *sliderElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*sliderElement)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
