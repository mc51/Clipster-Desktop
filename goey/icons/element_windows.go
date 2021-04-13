package icons

import (
	"github.com/lxn/win"
)

func (i *iconElement) SetOrder(hwnd win.HWND) win.HWND {
	return i.child.SetOrder(hwnd)
}
