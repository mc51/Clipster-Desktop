package goey

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"guitest/goey/base"
	"guitest/goey/dialog"
)

var (
	// ErrSetChildrenNotReentrant is returned if a reentrant call to the method
	// SetChild is called.
	ErrSetChildrenNotReentrant = errors.New("method SetChild is not reentrant")

	insideSetChildren uintptr
)

// Window represents a top-level window that contain other widgets.
type Window struct {
	windowImpl
}

// NewWindow create a new top-level window for the application.
func NewWindow(title string, child base.Widget) (*Window, error) {
	// Create the window
	w, err := newWindow(title, child)
	if err != nil {
		return nil, err
	}

	// The the default values for the horizontal and vertical scroll.
	// We want to do this before creating the child so that scrollbars can
	// be displayed (if necessary) with the relayout for the child.
	w.horizontalScroll, w.verticalScroll = scrollDefaults()

	// Mount the widget, and initialize its layout.
	if child != nil {
		newChild, err := child.Mount(w.control())
		if err != nil {
			w.Close()
			return nil, err
		}
		w.child = newChild
		w.setChildPost()
	}

	// Show the window
	w.show()

	if filename := os.Getenv("GOEY_SCREENSHOT"); filename != "" {
		asyncScreenshot(filename, w)
	}

	return w, nil
}

// NewWindow create a new hidden top-level window for the application.
// So we have loop running and Systray showing but no window
func NewHiddenWindow(title string, child base.Widget) (*Window, error) {
	// Create the window
	w, err := newWindow(title, child)
	if err != nil {
		return nil, err
	}

	// The the default values for the horizontal and vertical scroll.
	// We want to do this before creating the child so that scrollbars can
	// be displayed (if necessary) with the relayout for the child.
	w.horizontalScroll, w.verticalScroll = scrollDefaults()

	// Mount the widget, and initialize its layout.
	if child != nil {
		newChild, err := child.Mount(w.control())
		if err != nil {
			w.Close()
			return nil, err
		}
		w.child = newChild
		w.setChildPost()
	}

	// Show the window
	// w.show()

	if filename := os.Getenv("GOEY_SCREENSHOT"); filename != "" {
		asyncScreenshot(filename, w)
	}

	return w, nil
}

// Close destroys the window, and releases all associated resources.
func (w *Window) Close() {
	w.close()
}

// Child returns the mounted child for the window.  In general, this
// method should not be used.
func (w *Window) Child() base.Element {
	return w.child
}

// children assumes that the direct child of the window is a VBox, and then
// returns the children of that element.  It is used for testing.
func (w *Window) children() []base.Element {
	if w.child == nil {
		return nil
	}

	if vbox, ok := w.child.(*vboxElement); ok {
		return vbox.children
	}

	return nil
}

func (w *windowImpl) layoutChild(windowSize base.Size) base.Size {
	// Create the constraints
	constraints := base.Tight(windowSize)

	// Relax maximum size when scrolling is allowed
	if w.horizontalScroll {
		constraints.Max.Width = base.Inf
	}
	if w.verticalScroll {
		constraints.Max.Height = base.Inf
	}

	// Perform layout
	size := w.child.Layout(constraints)
	if !constraints.IsSatisfiedBy(size) {
		fmt.Println("constraints not satisfied,", constraints, ",", size)
	}
	return size
}

// Message returns a builder that can be used to construct a message
// dialog, and then show that dialog.
func (w *Window) Message(text string) *dialog.Message {
	ret := dialog.NewMessage(text)
	w.message(ret)
	return ret
}

// OpenFileDialog returns a builder that can be used to construct an open file
// dialog, and then show that dialog.
func (w *Window) OpenFileDialog() *dialog.OpenFile {
	ret := dialog.NewOpenFile()
	w.openfiledialog(ret)
	return ret
}

// SaveFileDialog returns a builder that can be used to construct a save file
// dialog, and then show that dialog.
func (w *Window) SaveFileDialog() *dialog.SaveFile {
	ret := dialog.NewSaveFile()
	w.savefiledialog(ret)
	return ret
}

// Scroll returns the flags that determine whether scrolling is allowed in the
// horizontal and vertical directions.
func (w *Window) Scroll() (horizontal, vertical bool) {
	return w.horizontalScroll, w.verticalScroll
}

func scrollDefaults() (horizontal, vertical bool) {
	env := os.Getenv("GOEY_SCROLL")
	if env == "" {
		return false, false
	}

	value, err := strconv.ParseUint(env, 10, 64)
	if err != nil || value >= 4 {
		return false, false
	}

	return (value & 2) == 2, (value & 1) == 1
}

// SetChild changes the child widget of the window.  As
// necessary, GUI widgets will be created or destroyed so that the GUI widgets
// match the widgets described by the parameter children.  The
// position of contained widgets will be updated to match the new layout
// properties.
func (w *Window) SetChild(child base.Widget) error {
	// One source of bugs in widgets is when the fire an event when being
	// updated.  This can lead to reentrant calls to SetChildren, typically
	// with incorrect information since the GUI is in an inconsistent state
	// when the event fires.  In short, this method is not reentrant.
	// The following will block changes to different windows, although
	// that shouldn't be susceptible to the same bugs.  Users in that
	// case should use Do to delay updates to other windows, but it shouldn't
	// happen in practice.
	if !atomic.CompareAndSwapUintptr(&insideSetChildren, 0, 1) {
		return ErrSetChildrenNotReentrant
	}
	defer func() {
		atomic.StoreUintptr(&insideSetChildren, 0)
	}()

	// Update the child element.
	newChild, err := base.DiffChild(w.control(), w.child, child)

	// Whether or not there has been an error, we need to run platform-specific
	// clean-up.  This is to recalculate min window size, update scrollbars, etc.
	w.child = newChild
	w.setChildPost()

	return err
}

// SetIcon changes the icon associated with the window.
//
// On Cocoa, individual windows do not have icons.  Instead, there is a single
// icon for the entire application.
func (w *Window) SetIcon(img image.Image) error {
	return w.setIcon(img)
}

// SetOnClosing changes the event callback for when the user tries to close the
// window.  This callback can also be used to save or close any resources
// before the window is closed.
//
// Returning true from the callback will prevent the window from closing.
func (w *Window) SetOnClosing(callback func() bool) {
	w.setOnClosing(callback)
}

// SetScroll sets whether scrolling is allowed in the horizontal and vertical directions.
func (w *Window) SetScroll(horizontal, vertical bool) {
	w.setScroll(horizontal, vertical)
}

// SetTitle changes the caption in the title bar for the window.
func (w *Window) SetTitle(title string) error {
	return w.setTitle(title)
}

// Title returns the current caption in the title bar for the window.
func (w *Window) Title() (string, error) {
	return w.title()
}

func sizeDefaults() (uint, uint) {
	const DEF_WIDTH = 640
	const DEF_HEIGHT = 480

	env := os.Getenv("GOEY_SIZE")
	if env == "" {
		return DEF_WIDTH, DEF_HEIGHT
	}

	parts := strings.Split(env, "x")
	if len(parts) != 2 {
		return DEF_WIDTH, DEF_HEIGHT
	}

	width, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return DEF_WIDTH, DEF_HEIGHT
	}

	height, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return DEF_WIDTH, DEF_HEIGHT
	}

	return uint(width), uint(height)
}
