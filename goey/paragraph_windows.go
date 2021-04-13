package goey

import (
	"syscall"
	"unsafe"

	"guitest/goey/base"
	win2 "guitest/goey/internal/syscall"
	"github.com/lxn/win"
)

func (w *P) calcStyle() uint32 {
	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.SS_LEFT)
	if w.Align == JustifyCenter {
		style = style | win.SS_CENTER
	} else if w.Align == JustifyRight {
		style = style | win.SS_RIGHT
	} else if w.Align == JustifyFull {
		style = style | win.SS_RIGHTJUST
	}
	return style
}

func (w *P) mount(parent base.Control) (base.Element, error) {
	// Create the control.
	hwnd, text, err := createControlWindow(0, &staticClassName[0], w.Text, w.calcStyle(), parent.HWnd)
	if err != nil {
		return nil, err
	}

	retval := &paragraphElement{Control: Control{hwnd}, text: text}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type paragraphElement struct {
	Control
	text []uint16
}

func (w *paragraphElement) measureReflowLimits() {
	hwnd := w.hWnd
	hdc := win.GetDC(hwnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	// Calculate the width of a single 'm' (find the em width)
	rect := win.RECT{0, 0, 0x7fffffff, 0x7fffffff}
	caption := [10]uint16{'m', 'm', 'm', 'm', 'm', 'm', 'm', 'm', 'm', 'm'}
	win.DrawTextEx(hdc, &caption[0], 10, &rect, win.DT_CALCRECT, nil)
	win.ReleaseDC(hwnd, hdc)
	paragraphMaxWidth = base.FromPixelsX(int(rect.Right)) * 8
}

func (w *paragraphElement) Props() base.Widget {
	align := JustifyLeft
	if style := win.GetWindowLong(w.hWnd, win.GWL_STYLE); style&win.SS_CENTER == win.SS_CENTER {
		align = JustifyCenter
	} else if style&win.SS_RIGHT == win.SS_RIGHT {
		align = JustifyRight
	} else if style&win.SS_RIGHTJUST == win.SS_RIGHTJUST {
		align = JustifyFull
	}

	return &P{
		Text:  w.Control.Text(),
		Align: align,
	}
}

func (w *paragraphElement) MinIntrinsicHeight(width base.Length) base.Length {
	if width == base.Inf {
		width = w.maxReflowWidth()
	}

	hdc := win.GetDC(w.hWnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	rect := win.RECT{0, 0, int32(width.PixelsX()), 0x7fffffff}
	win.DrawTextEx(hdc, &w.text[0], int32(len(w.text)), &rect, win.DT_CALCRECT|win.DT_WORDBREAK, nil)
	win.ReleaseDC(w.hWnd, hdc)

	return base.FromPixelsY(int(rect.Bottom))
}

func (w *paragraphElement) MinIntrinsicWidth(height base.Length) base.Length {
	if height != base.Inf {
		panic("not implemented")
	}

	width, _ := w.CalcRect(w.text)
	return min(base.FromPixelsX(int(width)), w.minReflowWidth())
}

func (w *paragraphElement) SetBounds(bounds base.Rectangle) {
	w.Control.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *paragraphElement) updateProps(data *P) error {
	text, err := syscall.UTF16FromString(data.Text)
	if err != nil {
		return err
	}
	w.text = text
	win2.SetWindowText(w.hWnd, &text[0])
	win.SetWindowLongPtr(w.hWnd, win.GWL_STYLE, uintptr(data.calcStyle()))

	return nil
}
