package main

import (
	"github.com/lxn/win"
)

func (w *minsizedboxElement) SetOrder(previous win.HWND) win.HWND {
	if w.child != nil {
		previous = w.child.SetOrder(previous)
	}
	return previous
}
