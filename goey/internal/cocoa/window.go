package cocoa

/*
#cgo CFLAGS: -x objective-c -DNTRACE
#cgo LDFLAGS: -framework Cocoa
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import (
	"image"
	"image/png"
	"os"
	"os/exec"
	"unsafe"
)

// Window is a wrapper for a NSWindow.
type Window struct {
	private int
}

type WindowCallbacks interface {
	OnShouldClose() bool
	OnWillClose()
	OnDidResize()
}

var (
	windowCallbacks = make(map[unsafe.Pointer]WindowCallbacks)
)

func NewWindow(title string, width, height uint) *Window {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.windowNew(ctitle, C.unsigned(width), C.unsigned(height))
	return (*Window)(handle)
}

func (w *Window) Close() {
	C.windowClose(unsafe.Pointer(w))
}

func (w *Window) ContentSize() (int, int) {
	size := C.windowContentSize(unsafe.Pointer(w))
	return int(size.width), int(size.height)
}

func (w *Window) ContentView() *View {
	return (*View)(C.windowContentView(unsafe.Pointer(w)))
}

func (w *Window) MakeFirstResponder(c *Control) {
	C.windowMakeFirstResponder(unsafe.Pointer(w), unsafe.Pointer(c))
}

func loadPNG(filename string) (image.Image, error) {
	file, err := os.Open("./ss.png")
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close of a file with only read permission.  Will not error.
		_ = file.Close()
	}()

	return png.Decode(file)
}

func (w *Window) Screenshot() image.Image {
	if ss := os.Getenv("SCREENSHOOTER"); ss != "" {
		cmd := exec.Command(ss, "-w", "-s", "./ss.png")
		err := cmd.Run()
		if err != nil {
			panic(err)
		}

		img, err := loadPNG("./ss.png")
		if err != nil {
			panic(err)
		}
		os.Remove("./ss.png")

		return img
	}

	size := C.windowFrameSize(unsafe.Pointer(w))

	img := image.NewRGBA(image.Rect(0, 0, int(size.width), int(size.height)))
	C.windowScreenshot(unsafe.Pointer(w), unsafe.Pointer(&img.Pix[0]), size.width, size.height)
	return img
}

func (w *Window) SetCallbacks(cb WindowCallbacks) {
	windowCallbacks[unsafe.Pointer(w)] = cb
}

func (w *Window) SetContentSize(width, height int) {
	C.windowSetContentSize(unsafe.Pointer(w), C.int(width), C.int(height))
}

func (w *Window) SetMinSize(width, height int) {
	C.windowSetMinSize(unsafe.Pointer(w), C.int(width), C.int(height))
}

func (w *Window) SetIcon(img image.Image) error {
	nsi, err := imageToNSImage(img)
	if err != nil {
		return err
	}

	C.windowSetIconImage(unsafe.Pointer(w), nsi)
	C.imageClose(nsi)
	return nil
}

func (w *Window) SetScrollVisible(horz, vert bool) {
	C.windowSetScrollVisible(unsafe.Pointer(w), toBool(horz), toBool(vert))
}

func (w *Window) SetTitle(title string) {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.windowSetTitle(unsafe.Pointer(w), ctitle)
}

func (w *Window) Title() string {
	return C.GoString(C.windowTitle(unsafe.Pointer(w)))
}

//export windowShouldClose
func windowShouldClose(handle unsafe.Pointer) bool {
	if cb := windowCallbacks[handle]; cb != nil {
		return cb.OnShouldClose()
	}

	return true
}

//export windowWillClose
func windowWillClose(handle unsafe.Pointer) {
	if cb := windowCallbacks[handle]; cb != nil {
		cb.OnWillClose()
	}
	delete(windowCallbacks, handle)
}

//export windowDidResize
func windowDidResize(handle unsafe.Pointer) {
	if cb := windowCallbacks[handle]; cb != nil {
		cb.OnDidResize()
	}
}
