package loop

import (
	"sync/atomic"
	"syscall"
	"testing"
	"unsafe"

	"clipster/goey/internal/nopanic"
	"github.com/lxn/win"
)

const (
	// Flag to control behaviour of UnlockOSThread in Run.
	unlockThreadAfterRun = true
)

var (
	atomPost win.ATOM
	hwndPost win.HWND
	namePost = [...]uint16{'G', 'o', 'e', 'y', 'P', 'o', 's', 't', 'W', 'i', 'n', 'd', 'o', 'w', 0}

	activeWindow uintptr
)

func initRun() error {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		// Not sure that the call to GetModuleHandle can ever fail when the
		// argument is nil.  The handle for the current .exe is certainly
		// valid?
		return syscall.GetLastError()
	}

	// Make sure that we have registered a class for the hidden window.
	if atomPost == 0 {
		wc := win.WNDCLASSEX{
			CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
			HInstance:     hInstance,
			LpfnWndProc:   syscall.NewCallback(postWindowProc),
			LpszClassName: &namePost[0],
		}

		atomPost = win.RegisterClassEx(&wc)
		if atomPost == 0 {
			return syscall.GetLastError()
		}
	}

	// Create the hidden window.
	hwndPost = win.CreateWindowEx(0, &namePost[0], nil, 0,
		win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT, win.CW_USEDEFAULT,
		win.HWND_DESKTOP, 0, hInstance, nil)
	if hwndPost == 0 {
		return syscall.GetLastError()
	}
	return nil
}

func terminateRun() {
	win.DestroyWindow(hwndPost)
	hwndPost = 0
}

func run() {
	// Run the message loop
	for loop() {
	}
}

func runTesting(func() error) error {
	panic("unreachable")
}

func do(action func() error) error {
	// Marshal the action to the GUI thread, and collect the return value.
	err := make(chan error, 1)
	win.PostMessage(hwndPost, win.WM_USER, uintptr(unsafe.Pointer(&action)), uintptr(unsafe.Pointer(&err)))
	return nopanic.Unwrap(<-err)
}

func loop() (ok bool) {
	// Obtain a copy of the next message from the queue.
	var msg win.MSG
	win.GetMessage(&msg, 0, 0, 0)

	// Processing for application wide messages are handled in this block.
	if msg.Message == win.WM_QUIT {
		return false
	}

	// Dispatch message.
	if !win.IsDialogMessage(win.HWND(activeWindow), &msg) {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
	return true
}

func stop() {
	win.PostQuitMessage(0)
}

func SetActiveWindow(hwnd win.HWND) {
	atomic.StoreUintptr(&activeWindow, uintptr(hwnd))
}

func postWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case win.WM_USER:
		err := nopanic.Wrap(*(*func() error)(unsafe.Pointer(wParam)))
		(*(*chan error)(unsafe.Pointer(lParam))) <- err
		return 0
	}

	// Let the default window proc handle all other messages
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func testMain(m *testing.M) int {
	// On Windows, we need to be locked to a thread, but not to a particular
	// thread.  No need for special coordination.
	return m.Run()
}
