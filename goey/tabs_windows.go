package goey

import (
	"image/color"
	"syscall"
	"unsafe"

	"clipster/goey/base"
	"github.com/lxn/win"
)

var (
	tabs struct {
		className     []uint16
		oldWindowProc uintptr
		hbrush        win.HBRUSH
		hbrushFlag    bool
	}
)

func init() {
	tabs.className = []uint16{'S', 'y', 's', 'T', 'a', 'b', 'C', 'o', 'n', 't', 'r', 'o', 'l', '3', '2', 0}
}

func (w *Tabs) mount(parent base.Control) (base.Element, error) {
	// Create the control
	const STYLE = win.WS_CHILD | win.WS_VISIBLE
	hwnd, _, err := createControlWindow(win.WS_EX_CONTROLPARENT, &tabs.className[0], "", STYLE, parent.HWnd)
	if err != nil {
		return nil, err
	}

	for i, v := range w.Children {
		text, err := syscall.UTF16PtrFromString(v.Caption)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		item := win.TCITEM{
			Mask:    win.TCIF_TEXT,
			PszText: text,
			LParam:  uintptr(i),
		}
		win.SendMessage(hwnd, win.TCM_INSERTITEM, uintptr(i), uintptr(unsafe.Pointer(&item)))
	}
	if w.Value > 0 {
		win.SendMessage(hwnd, win.TCM_SETCURSEL, uintptr(w.Value), 0)
	}

	child := base.Element(nil)
	if len(w.Children) > 0 {
		var err error
		if w.Value >= 0 {
			child, err = base.Mount(base.Control{hwnd}, w.Children[w.Value].Child)
		} else {
			child, err = base.Mount(base.Control{hwnd}, w.Children[0].Child)
		}
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}
	}

	retval := &tabsElement{
		Control:  Control{hwnd},
		child:    child,
		parent:   parent,
		value:    w.Value,
		insets:   w.Insets,
		widgets:  w.Children,
		onChange: w.OnChange,
	}

	// Subclass the window procedure
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))
	subclassWindowProcedure(hwnd, &tabs.oldWindowProc, tabsWindowProc)

	return retval, nil
}

type tabsElement struct {
	Control
	child    base.Element
	parent   base.Control
	value    int
	insets   Insets
	widgets  []TabItem
	onChange func(int)

	cachedInsets base.Point
	cachedBounds base.Rectangle
	hbrush       win.HBRUSH
}

func (w *tabsElement) contentInsets() base.Point {
	if w.cachedInsets.Y == 0 {
		rect := win.RECT{}

		win.SendMessage(w.hWnd, win.TCM_ADJUSTRECT, win.TRUE, uintptr(unsafe.Pointer(&rect)))
		w.cachedInsets = base.Point{
			X: base.FromPixelsX(int(rect.Right - rect.Left)),
			Y: base.FromPixelsY(int(rect.Bottom - rect.Top)),
		}
	}

	return w.cachedInsets
}

func (w *tabsElement) controlTabsMinWidth() base.Length {
	// No API to get this information has been found.
	return 75 * DIP
}

func (w *tabsElement) Props() base.Widget {
	count := win.SendMessage(w.hWnd, win.TCM_GETITEMCOUNT, 0, 0)
	children := make([]TabItem, count)
	for i := uintptr(0); i < count; i++ {
		text := [128]uint16{}
		item := win.TCITEM{
			Mask:       win.TCIF_TEXT,
			PszText:    &text[0],
			CchTextMax: 128,
		}

		win.SendMessage(w.hWnd, win.TCM_GETITEM, i, uintptr(unsafe.Pointer(&item)))
		children[i].Caption = syscall.UTF16ToString(text[:])
		children[i].Child = w.widgets[i].Child
	}

	return &Tabs{
		Value:    int(win.SendMessage(w.hWnd, win.TCM_GETCURSEL, 0, 0)),
		Children: children,
		OnChange: w.onChange,
	}
}

func (w *tabsElement) SetOrder(previous win.HWND) win.HWND {
	previous = w.Control.SetOrder(previous)
	if w.child != nil {
		previous = w.child.SetOrder(previous)
	}
	return previous
}

func (w *tabsElement) SetBounds(bounds base.Rectangle) {
	w.Control.SetBounds(bounds)
	if w.hbrush != 0 {
		win.DeleteObject(win.HGDIOBJ(w.hbrush))
		w.hbrush = 0
	}

	if w.child != nil {
		// Determine the bounds for the child widget
		rect := win.RECT{}
		win.SendMessage(w.hWnd, win.TCM_ADJUSTRECT, win.FALSE, uintptr(unsafe.Pointer(&rect)))
		w.cachedBounds = base.Rectangle{
			Min: bounds.Min.Add(base.Point{base.FromPixelsX(int(rect.Left)), base.FromPixelsY(int(rect.Top))}),
			Max: bounds.Max.Add(base.Point{base.FromPixelsX(int(rect.Right)), base.FromPixelsY(int(rect.Bottom))}),
		}
		// Offset to handle insets
		w.cachedBounds.Min.X += w.insets.Left - bounds.Min.X
		w.cachedBounds.Min.Y += w.insets.Top - bounds.Min.Y
		w.cachedBounds.Max.X -= w.insets.Right + bounds.Min.X
		w.cachedBounds.Max.Y -= w.insets.Bottom + bounds.Min.Y

		// Update bounds for the child
		w.child.SetBounds(w.cachedBounds)
	}
}

func (w *tabsElement) updateChildren(children []TabItem) error {
	len1 := len(w.widgets)
	len2 := len(children)

	if len1 <= len2 {
		// Change caption for tabs that already exist
		for i, v := range children[:len1] {
			text, err := syscall.UTF16PtrFromString(v.Caption)
			if err != nil {
				return err
			}

			item := win.TCITEM{
				Mask:    win.TCIF_TEXT,
				PszText: text,
			}
			win.SendMessage(w.hWnd, win.TCM_SETITEM, uintptr(i), uintptr(unsafe.Pointer(&item)))
		}

		// Add new tabs to extend the list
		for i, v := range children[len1:] {
			text, err := syscall.UTF16PtrFromString(v.Caption)
			if err != nil {
				return err
			}

			item := win.TCITEM{
				Mask:    win.TCIF_TEXT,
				PszText: text,
			}
			win.SendMessage(w.hWnd, win.TCM_INSERTITEM, uintptr(i+len1), uintptr(unsafe.Pointer(&item)))
		}
	} else {
		// Change caption for tabs that already exist
		for i, v := range children {
			text, err := syscall.UTF16PtrFromString(v.Caption)
			if err != nil {
				return err
			}

			item := win.TCITEM{
				Mask:    win.TCIF_TEXT,
				PszText: text,
			}
			win.SendMessage(w.hWnd, win.TCM_SETITEM, uintptr(i), uintptr(unsafe.Pointer(&item)))
		}

		// Delete excess tabs.
		for i := len2; i < len1; i++ {
			win.SendMessage(w.hWnd, win.TCM_DELETEITEM, uintptr(i), 0)
		}
	}

	w.widgets = children
	return nil
}

func (w *tabsElement) updateProps(data *Tabs) error {
	// Update the tabs
	err := w.updateChildren(data.Children)
	if err != nil {
		return err
	}

	// Update which tab is currently selected
	if data.Value >= 0 && w.value != data.Value {
		win.SendMessage(w.hWnd, win.TCM_SETCURSEL, uintptr(data.Value), 0)
		w.value = data.Value
	}

	// Update event handlers
	w.onChange = data.OnChange

	return nil
}

func tabsBackgroundBrush(hwnd win.HWND, hdc win.HDC) (win.HBRUSH, bool, error) {
	// If there is a global brush that can be used for all tab controls,
	// use that brush
	if tabs.hbrush != 0 {
		return tabs.hbrush, true, nil
	}

	// We need to print the client area for the tab control to a bitmap, which
	// can be used to create a brush.
	cr := win.RECT{}
	win.GetClientRect(hwnd, &cr)

	// Configure a device context to capture contents of the control
	cdc := win.CreateCompatibleDC(hdc)
	if cdc == 0 {
		return 0, false, syscall.GetLastError()
	}
	defer func() {
		win.DeleteDC(cdc)
	}()
	hbitmap := win.CreateCompatibleBitmap(hdc, cr.Right-cr.Left, cr.Bottom-cr.Top)
	if hbitmap == 0 {
		return 0, false, syscall.GetLastError()
	}
	defer func() {
		win.DeleteObject(win.HGDIOBJ(hbitmap))
	}()
	win.SelectObject(cdc, win.HGDIOBJ(hbitmap))

	// Get bitmap of the control
	win.SendMessage(hwnd, win.WM_PRINTCLIENT, uintptr(cdc), win.PRF_CLIENT)

	// If possible, better to use a solid brush that does not need to be
	// generated everytime the size of the control is changed.
	if !tabs.hbrushFlag {
		tabs.hbrushFlag = true

		// Is the current bitmap a constant color in the client area
		win.SendMessage(hwnd, win.TCM_ADJUSTRECT, win.FALSE, uintptr(unsafe.Pointer(&cr)))
		if clr := win.GetPixel(cdc, cr.Left, cr.Top); clr == 0 {
			// Don't believe it.  Windows lies.
			// We are running on windows without themes enabled.
			tabs.hbrush = win.GetSysColorBrush(win.COLOR_3DFACE)
			return tabs.hbrush, true, nil
		} else if clr == win.GetPixel(cdc, (cr.Left+cr.Right)/2, cr.Top) && clr == win.GetPixel(cdc, cr.Left, (cr.Top+cr.Bottom)/2) {
			hbrush := createBrush(color.RGBA{
				R: uint8(clr & 0xFF),
				G: uint8((clr >> 8) & 0xFF),
				B: uint8((clr >> 16) & 0xFF),
				A: 0xFF,
			})
			tabs.hbrush = hbrush
			return hbrush, true, nil
		}
	}

	// Convert the bitmap to a brush
	brush := win.CreatePatternBrush(hbitmap)
	if brush == 0 {
		return 0, false, syscall.GetLastError()
	}

	return brush, false, nil
}

func tabsWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if brush := tabsGetPtr(hwnd).hbrush; brush != 0 {
			win.DeleteObject(win.HGDIOBJ(brush))

		}
		tabsGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_CTLCOLORSTATIC:
		win.SetBkMode(win.HDC(wParam), win.TRANSPARENT)
		w := tabsGetPtr(hwnd)
		if w.hbrush == 0 {
			if hbrush, solid, err := tabsBackgroundBrush(hwnd, win.HDC(wParam)); err != nil {
				panic(err)
			} else if solid {
				return uintptr(hbrush)
			} else {
				w.hbrush = hbrush
			}
		}
		// Set offset for the brush
		cr := win.RECT{}
		win.GetWindowRect(win.HWND(lParam), &cr)
		origin := win.POINT{cr.Left, cr.Top}
		win.ScreenToClient(hwnd, &origin)
		win.SetBrushOrgEx(win.HDC(wParam), -origin.X, -origin.Y, nil)
		return uintptr(w.hbrush)

	case win.WM_HSCROLL:
		if lParam != 0 {
			// Message was sent by a child window.  As for all other controls
			// that notify the parent, resend to the child with the expectation
			// that the child has been subclassed.
			return win.SendMessage(win.HWND(lParam), win.WM_HSCROLL, wParam, 0)
		}
		// Defer to default window proc

	case win.WM_COMMAND:
		return windowprocWmCommand(wParam, lParam)

	case win.WM_NOTIFY:
		if n := (*win.NMHDR)(unsafe.Pointer(lParam)); true {
			if n.HwndFrom == hwnd {
				if n.Code == uint32(0x100000000+win.TCN_SELCHANGE) {
					cursel := int(win.SendMessage(hwnd, win.TCM_GETCURSEL, 0, 0))
					if w := tabsGetPtr(hwnd); w.value != cursel {
						if w.onChange != nil {
							w.onChange(cursel)
						}
						if w.value != cursel {
							child, err := base.DiffChild(w.parent, w.child, w.widgets[cursel].Child)
							if err != nil {
								panic("Unhandled error!")
							}
							if child != nil {
								child.SetOrder(w.hWnd)
								child.Layout(base.Tight(base.Size{
									Width:  w.cachedBounds.Dx(),
									Height: w.cachedBounds.Dy(),
								}))
								child.SetBounds(w.cachedBounds)
								win.InvalidateRect(win.GetParent(w.hWnd), nil, false)
							}
							w.child = child
							w.value = cursel
						}
					}
				}
			} else {
				return windowprocWmNotify(wParam, lParam)
			}
		}
		return 0
	}

	return win.CallWindowProc(tabs.oldWindowProc, hwnd, msg, wParam, lParam)
}

func tabsGetPtr(hwnd win.HWND) *tabsElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*tabsElement)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
