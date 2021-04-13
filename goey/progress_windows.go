package goey

import (
	"unsafe"

	"guitest/goey/base"
	"github.com/lxn/win"
)

var (
	progress struct {
		className     []uint16
		oldWindowProc uintptr
	}
)

func init() {
	progress.className = []uint16{'m', 's', 'c', 't', 'l', 's', '_', 'p', 'r', 'o', 'g', 'r', 'e', 's', 's', '3', '2', 0}
}

func (w *Progress) mount(parent base.Control) (base.Element, error) {
	// Create the control.
	style := uint32(win.WS_CHILD | win.WS_VISIBLE)
	hwnd, _, err := createControlWindow(0, &progress.className[0], "", style, parent.HWnd)
	if err != nil {
		return nil, err
	}
	win.SendMessage(hwnd, win.PBM_SETRANGE32, uintptr(w.Min), uintptr(w.Max))
	win.SendMessage(hwnd, win.PBM_SETPOS, uintptr(w.Value), 0)

	retval := &progressElement{
		Control: Control{hwnd},
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type progressElement struct {
	Control
}

func (w *progressElement) Layout(bc base.Constraints) base.Size {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	width := w.MinIntrinsicWidth(0)
	if bc.Max.Width > 355*DIP {
		width = 355 * DIP
	}
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *progressElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 15 * DIP
}

func (w *progressElement) MinIntrinsicWidth(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 160 * DIP
}

func (w *progressElement) Props() base.Widget {
	min := win.SendMessage(w.hWnd, win.PBM_GETRANGE, win.TRUE, 0)
	max := win.SendMessage(w.hWnd, win.PBM_GETRANGE, win.FALSE, 0)
	value := win.SendMessage(w.hWnd, win.PBM_GETPOS, 0, 0)

	return &Progress{
		Value: int(value),
		Min:   int(min),
		Max:   int(max),
	}
}

func (w *progressElement) updateProps(data *Progress) error {
	win.SendMessage(w.hWnd, win.PBM_SETRANGE32, uintptr(data.Min), uintptr(data.Max))
	win.SendMessage(w.hWnd, win.PBM_SETPOS, uintptr(data.Value), 0)
	return nil
}
