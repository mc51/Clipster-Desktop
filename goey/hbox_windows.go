package goey

import (
	"github.com/lxn/win"
)

func (w *hboxElement) SetOrder(previous win.HWND) win.HWND {
	for _, v := range w.children {
		previous = v.SetOrder(previous)
	}
	return previous
}
