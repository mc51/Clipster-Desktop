package goey

import (
	"syscall"
	"unsafe"

	"guitest/goey/base"
	"github.com/lxn/win"
)

func (w *TextArea) mount(parent base.Control) (base.Element, error) {
	// Create the control, and set properties
	hwnd, _, err := createControlWindow(win.WS_EX_CLIENTEDGE, &edit.className[0], w.Value, w.style(), parent.HWnd)
	if err != nil {
		return nil, err
	}
	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Create placeholder, if required.
	if w.Placeholder != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(w.Placeholder)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		win.SendMessage(hwnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	}

	// Create the return value.
	retval := &textareaElement{textinputElementBase{
		Control:  Control{hwnd},
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	},
		minlinesDefault(w.MinLines),
	}

	// Link the control back to Go for event handling
	subclassWindowProcedure(hwnd, &edit.oldWindowProc, textinputWindowProc)
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

func (w *TextArea) style() uint32 {
	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_MULTILINE | win.ES_WANTRETURN | win.ES_AUTOVSCROLL)
	if w.ReadOnly {
		style = style | win.ES_READONLY
	}
	return style
}

type textareaElement struct {
	textinputElementBase
	minLines int
}

func (w *textareaElement) Layout(bc base.Constraints) base.Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(base.Size{width, height})
}

func (w *textareaElement) MinIntrinsicHeight(base.Length) base.Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	const lineHeight = 16 * DIP
	return 23*DIP + lineHeight.Scale(w.minLines-1, 1)
}

func (w *textareaElement) Props() base.Widget {
	var buffer [80]uint16
	win.SendMessage(w.hWnd, win.EM_GETCUEBANNER, uintptr(unsafe.Pointer(&buffer[0])), 80)
	ndx := 0
	for i, v := range buffer {
		if v == 0 {
			ndx = i
			break
		}
	}
	placeholder := syscall.UTF16ToString(buffer[:ndx])

	return &TextArea{
		Value:       w.Control.Text(),
		Placeholder: placeholder,
		Disabled:    !win.IsWindowEnabled(w.hWnd),
		ReadOnly:    (win.GetWindowLong(w.hWnd, win.GWL_STYLE) & win.ES_READONLY) != 0,
		MinLines:    w.minLines,
		OnChange:    w.onChange,
		OnFocus:     w.onFocus,
		OnBlur:      w.onBlur,
	}
}

func (w *textareaElement) updateProps(data *TextArea) error {
	if data.Value != w.Text() {
		w.SetText(data.Value)
	}
	err := updatePlaceholder(w.hWnd, data.Placeholder)
	if err != nil {
		return err
	}
	w.SetDisabled(data.Disabled)
	win.SendMessage(w.hWnd, win.EM_SETREADONLY, uintptr(win.BoolToBOOL(data.ReadOnly)), 0)

	w.minLines = minlinesDefault(data.MinLines)
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}
