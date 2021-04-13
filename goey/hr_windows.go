package goey

import (
	"syscall"
	"unsafe"

	"guitest/goey/base"
	"github.com/lxn/win"
)

var (
	hr struct {
		className []uint16
		atom      win.ATOM
	}
)

func init() {
	hr.className = []uint16{'G', 'o', 'e', 'y', 'H', 'R', 0}
}

func registerHRClass(hInst win.HINSTANCE, wndproc uintptr) (win.ATOM, error) {
	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.HInstance = hInst
	wc.LpfnWndProc = wndproc
	wc.HCursor = win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	wc.HbrBackground = (win.HBRUSH)(win.GetStockObject(win.NULL_BRUSH))
	wc.LpszClassName = &hr.className[0]

	atom := win.RegisterClassEx(&wc)
	if atom == 0 {
		return 0, syscall.GetLastError()
	}
	return atom, nil
}

func (w *HR) mount(parent base.Control) (base.Element, error) {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		return nil, syscall.GetLastError()
	}
	if hr.atom == 0 {
		atom, err := registerHRClass(hInstance, syscall.NewCallback(hrWindowProc))
		if err != nil {
			return nil, err
		}
		hr.atom = atom
	}

	hwnd := win.CreateWindowEx(0, &hr.className[0], nil, win.WS_CHILD|win.WS_VISIBLE,
		10, 10, 100, 100,
		parent.HWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}

	retval := &hrElement{Control: Control{hwnd}}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type hrElement struct {
	Control
}

func hrWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*hrElement)(unsafe.Pointer(w))
			ptr.hWnd = 0
		}
		// Defer to the old window proc

	case win.WM_PAINT:
		ps := win.PAINTSTRUCT{}
		rect := win.RECT{}
		hdc := win.BeginPaint(hwnd, &ps)
		win.SetBkMode(hdc, win.TRANSPARENT)
		win.GetClientRect(hwnd, &rect)
		win.MoveToEx(hdc, int(rect.Left), int(rect.Top+rect.Bottom)/2, nil)
		win.LineTo(hdc, rect.Right, (rect.Top+rect.Bottom)/2)
		win.EndPaint(hwnd, &ps)
	}

	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}
