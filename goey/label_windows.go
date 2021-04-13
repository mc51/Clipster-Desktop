package goey

import (
	"syscall"
	"unsafe"

	"guitest/goey/base"
	win2 "guitest/goey/internal/syscall"
	"github.com/lxn/win"
)

var (
	staticClassName []uint16
)

func init() {
	staticClassName = []uint16{'S', 'T', 'A', 'T', 'I', 'C', 0}
}

func (w *Label) mount(parent base.Control) (base.Element, error) {
	// Create the control
	const STYLE = win.WS_CHILD | win.WS_VISIBLE | win.SS_LEFT
	hwnd, text, err := createControlWindow(0, &staticClassName[0], w.Text, STYLE, parent.HWnd)
	if err != nil {
		return nil, err
	}

	retval := &labelElement{Control: Control{hwnd}, text: text}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type labelElement struct {
	Control
	text []uint16
}

func (w *labelElement) Props() base.Widget {
	return &Label{
		Text: w.Control.Text(),
	}
}

func (w *labelElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *labelElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13 * DIP
}

func (w *labelElement) MinIntrinsicWidth(base.Length) base.Length {
	width, _ := w.CalcRect(w.text)
	return base.FromPixelsX(int(width))
}

func (w *labelElement) SetBounds(bounds base.Rectangle) {
	// Because of descenders in text, we may want to increase the height
	// of the label.
	_, height := w.CalcRect(w.text)
	if h := base.FromPixelsY(int(height)); h > bounds.Dy() {
		bounds.Max.Y = bounds.Min.Y + h
	}

	w.Control.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *labelElement) updateProps(data *Label) error {
	text, err := syscall.UTF16FromString(data.Text)
	if err != nil {
		return err
	}
	w.text = text
	win2.SetWindowText(w.hWnd, &text[0])
	// TODO:  Update alignment

	return nil
}
