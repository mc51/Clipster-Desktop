package goey

import (
	"guitest/goey/base"
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	button struct {
		className     []uint16
		oldWindowProc uintptr
	}
)

func init() {
	button.className = []uint16{'B', 'U', 'T', 'T', 'O', 'N', 0}
}

func buttonStyle(isDefault bool) uint32 {
	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.BS_PUSHBUTTON | win.BS_TEXT | win.BS_NOTIFY)
	if isDefault {
		style = style | win.BS_DEFPUSHBUTTON
	}
	return style
}

func (w *Button) mount(parent base.Control) (base.Element, error) {
	// Create the control.
	hwnd, text, err := createControlWindow(0, &button.className[0], w.Text, buttonStyle(w.Default), parent.HWnd)
	if err != nil {
		return nil, err
	}
	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &button.oldWindowProc, buttonWindowProc)

	retval := &buttonElement{
		Control: Control{hwnd},
		text:    text,
		onClick: w.OnClick,
		onFocus: w.OnFocus,
		onBlur:  w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type buttonElement struct {
	Control
	text []uint16

	onClick func()
	onFocus func()
	onBlur  func()
}

func (w *buttonElement) Click() {
	win.SendMessage(w.hWnd, win.BM_CLICK, 0, 0)
}

func (w *buttonElement) Props() base.Widget {
	return &Button{
		Text:     w.Control.Text(),
		Disabled: !win.IsWindowEnabled(w.hWnd),
		Default:  (win.GetWindowLong(w.hWnd, win.GWL_STYLE) & win.BS_DEFPUSHBUTTON) != 0,
		OnClick:  w.onClick,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *buttonElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *buttonElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 * DIP
}

func (w *buttonElement) MinIntrinsicWidth(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	width, _ := w.CalcRect(w.text)
	return max(
		75*DIP,
		base.FromPixelsX(int(width)+7),
	)
}

func (w *buttonElement) updateProps(data *Button) error {
	text, err := syscall.UTF16FromString(data.Text)
	if err != nil {
		return err
	}

	w.SetText(data.Text)
	w.text = text
	w.SetDisabled(data.Disabled)
	win.SendMessage(w.hWnd, win.BM_SETSTYLE, uintptr(buttonStyle(data.Default)), win.TRUE)
	w.onClick = data.OnClick
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}

func buttonWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		buttonGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := buttonGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := buttonGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		// WM_COMMAND is sent to the parent, which will only forward certain
		// message.  This code should only ever see EN_UPDATE, but we will
		// still check.
		switch notification := win.HIWORD(uint32(wParam)); notification {
		case win.BN_CLICKED:
			if w := buttonGetPtr(hwnd); w.onClick != nil {
				w.onClick()
			}
		}
		return 0
	}

	return win.CallWindowProc(button.oldWindowProc, hwnd, msg, wParam, lParam)
}

func buttonGetPtr(hwnd win.HWND) *buttonElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*buttonElement)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
