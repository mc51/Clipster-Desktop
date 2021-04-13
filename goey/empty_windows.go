package goey

import (
	"github.com/lxn/win"
)

func (w *emptyElement) SetOrder(hwnd win.HWND) win.HWND {
	return hwnd
}
