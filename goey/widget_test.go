package goey

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"clipster/goey/base"
	"clipster/goey/loop"
)

type Proper interface {
	Props() base.Widget
}

type Clickable interface {
	Click()
}

type Focusable interface {
	TakeFocus() bool
}

type Typeable interface {
	TakeFocus() bool
	TypeKeys(text string) chan error
}

func normalize(t *testing.T, rhs base.Widget) {
	if base.PLATFORM == "windows" {
		// On windows, the message EM_GETCUEBANNER does not work unless the manifest
		// is set correctly.  This cannot be done for the package, since that
		// manifest will conflict with the manifest of any app.
		if value := reflect.ValueOf(rhs).Elem().FieldByName("Placeholder"); value.IsValid() {
			placeholder := value.String()
			if placeholder != "" {
				t.Logf("Zeroing 'Placeholder' field during test")
			}
			value.SetString("")
		}
	} else if base.PLATFORM == "gtk" {
		// With GTK, this package is using a GtkTextView to create
		// the multi-line text editor, and that widget does not support
		// placeholders.
		if elem, ok := rhs.(*TextArea); ok {
			if elem.Placeholder != "" {
				t.Logf("Zeroing 'Placeholder' field during test")
			}
			elem.Placeholder = ""
		}
	}

	if base.PLATFORM != "cocoa" {
		// On both windows and GTK, the props method only return RGBA images.
		if value := reflect.ValueOf(rhs).Elem().FieldByName("Image"); value.IsValid() {
			if prop, ok := value.Interface().(*image.Gray); ok {
				t.Logf("Converting 'Image' field from *image.Gray to *image.RGBA")
				bounds := prop.Bounds()
				img := image.NewRGBA(bounds)
				draw.Draw(img, bounds, prop, bounds.Min, draw.Src)
				value.Set(reflect.ValueOf(img))
			}
		}
	}

	if value := reflect.ValueOf(rhs).Elem().FieldByName("Child"); value.IsValid() {
		if child := value.Interface(); child != nil {
			normalize(t, child.(base.Widget))
		}
	}
}

func equal(t *testing.T, lhs, rhs base.Widget) bool {
	// Normalize (or canonicalize) the props used to construct the element.
	normalize(t, rhs)
	// Compare the widgets' properties.
	return reflect.DeepEqual(lhs, rhs)
}

func testingMountWidgets(t *testing.T, widgets ...base.Widget) {
	init := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
			return nil
		}
		if window == nil {
			t.Errorf("Unexpected nil for window")
			return nil
		}

		// Check that the controls that were mounted match with the list
		if children := window.children(); len(children) != len(widgets) {
			t.Errorf("Wanted len(children) == len(widgets), got %d and %d", len(children), len(widgets))
		} else {
			for i := range children {
				if n1, n2 := children[i].Kind(), widgets[i].Kind(); n1 != n2 {
					t.Errorf("Wanted children[%d].Kind() == widgets[%d].Kind(), got %s, want %s", i, i, n1, n2)
				} else if widget, ok := children[i].(Proper); ok {
					data := widget.Props()
					if n1, n2 := data.Kind(), widgets[i].Kind(); n1 != n2 {
						t.Errorf("Wanted data.Kind() == widgets[%d].Kind(), got %s, want %s", i, n1, n2)
					}
					if !equal(t, data, widgets[i]) {
						t.Errorf("Wanted data == widgets[%d], got %v, want %v", i, data, widgets[i])
					}
				} else {
					t.Logf("Cannot verify props of child")
				}
			}
		}
		go func(window *Window) {
			if testing.Verbose() && !testing.Short() {
				time.Sleep(25 * time.Millisecond)
			}
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
}

func testingMountWidget(t *testing.T, widget base.Widget) (ok bool) {
	init := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		window, err := NewWindow(t.Name(), &VBox{Children: []base.Widget{widget}})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
			return nil
		}
		if window == nil {
			t.Errorf("Unexpected nil for window")
			return nil
		}

		// Check that the controls that were mounted match with the list
		if children := window.children(); len(children) != 1 {
			t.Errorf("Wanted len(children) == 1, got %d", len(children))
		} else {
			ok = children[0].Kind() == widget.Kind() &&
				equal(t, children[0].(Proper).Props(), widget)
		}

		go func(window *Window) {
			if testing.Verbose() {
				time.Sleep(25 * time.Millisecond)
			}
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	return /* naked return, code set in init callback */
}

func testingMountWidgetsFail(t *testing.T, outError error, widgets ...base.Widget) {
	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if window != nil {
			t.Errorf("Unexpected non-nil window")
		}
		if err != outError {
			if err == nil {
				t.Errorf("Unexpected nil error, want %s", outError)
			} else {
				t.Errorf("Unexpected error, want %v, got %s", outError, err)
			}
			return nil
		}
		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
}

func testingCloseWidgets(t *testing.T, widgets ...base.Widget) {
	init := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
			return nil
		}
		if window == nil {
			t.Errorf("Unexpected nil for window")
			return nil
		}

		// Check that the controls that were mounted match with the list
		if len(window.children()) != len(widgets) {
			t.Errorf("Want len(window.Children())!=nil")
		}

		err = window.SetChild(&VBox{Children: nil})
		if err != nil {
			t.Errorf("Failed to set children, %s", err)
			return nil
		}
		if len(window.children()) != 0 {
			t.Errorf("Want len(window.Children())!=0")
		}

		go func(window *Window) {
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
}

func testingCheckFocusAndBlur(t *testing.T, widgets ...base.Widget) {
	log := bytes.NewBuffer(nil)
	skipFlag := false

	for i := byte(0); i < 3; i++ {
		s := reflect.ValueOf(widgets[i])
		letter := 'a' + i
		s.Elem().FieldByName("OnFocus").Set(reflect.ValueOf(func() {
			log.Write([]byte{'f', letter})
		}))
		s.Elem().FieldByName("OnBlur").Set(reflect.ValueOf(func() {
			log.Write([]byte{'b', letter})
		}))
	}

	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}

		go func(window *Window) {
			// Wait for the window to be active.
			// This does not appear to be necessary on WIN32.  With GTK, the
			// window needs time to display before it will respond properly to
			// focus events.
			time.Sleep(20 * time.Millisecond)

			// Run the actions, which are counted.
			for i := 0; i < 3; i++ {
				err := loop.Do(func() error {
					// Find the child element to be focused
					child := window.child.(*vboxElement).children[i]
					if elem, ok := child.(Focusable); ok {
						ok := elem.TakeFocus()
						if !ok {
							t.Errorf("Failed to set focus on the control")
						}
					} else {
						skipFlag = true
					}
					return nil
				})
				if err != nil {
					t.Errorf("Error in Do, %s", err)
				}
			}

			// Close the window
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if skipFlag {
		t.Skip("Control does not support TakeFocus")
	}

	const want = "fabafbbbfcbc"
	if s := log.String(); s != want {
		t.Errorf("Incorrect log string, want %s, got log==%s", want, s)
	}
}

func testingTypeKeys(t *testing.T, text string, widget base.Widget) {
	skipFlag := false

	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: []base.Widget{widget}})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}

		var typingErr chan error
		go func(window *Window) {
			// On WIN32, let the window complete any opening animation.
			time.Sleep(20 * time.Millisecond)

			err := loop.Do(func() error {
				child := window.child.(*vboxElement).children[0]
				if elem, ok := child.(Typeable); ok {
					if elem.TakeFocus() {
						typingErr = elem.TypeKeys(text)
					} else {
						t.Errorf("Control failed to take focus.")
					}
				} else {
					skipFlag = true
				}
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}

			// Wait for typing to complete, and check for errors
			if typingErr != nil {
				for v := range typingErr {
					t.Errorf("Failed to type keys on the control, %v", v)
				}
			}

			// Close the window
			err = loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if skipFlag {
		t.Skip("Control does not support TypeKeys")
	}
}

func testingCheckClick(t *testing.T, widgets ...base.Widget) {
	log := bytes.NewBuffer(nil)
	skipFlag := false

	for i := byte(0); i < 3; i++ {
		letter := 'a' + i
		if elem, ok := widgets[i].(*Checkbox); ok {
			// Chain the onclick callback for the element.
			chainCallback := elem.OnChange
			// Add wrapper to write to the test log.
			elem.OnChange = func(value bool) {
				log.Write([]byte{'c', letter})
				chainCallback(value)
			}
		} else {
			s := reflect.ValueOf(widgets[i])
			s.Elem().FieldByName("OnClick").Set(reflect.ValueOf(func() {
				log.Write([]byte{'c', letter})
			}))
		}
	}

	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}

		go func(window *Window) {
			// Run the actions, which are counted.
			for i := 0; i < 3; i++ {
				err := loop.Do(func() error {
					// Find the child element to be clicked
					child := window.child.(*vboxElement).children[i]
					if elem, ok := child.(Clickable); ok {
						elem.Click()
					} else {
						skipFlag = true
					}
					return nil
				})
				if err != nil {
					t.Errorf("Error in Do, %s", err)
				}
			}

			// Close the window
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if skipFlag {
		t.Skip("Control does not support Click")
	}

	const want = "cacbcc"
	if s := log.String(); s != want {
		t.Errorf("Incorrect log string, want %s, got log==%s", want, s)
	}
}

func testingUpdateWidgets(t *testing.T, widgets []base.Widget, update []base.Widget) {
	init := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
			return nil
		}
		if window == nil {
			t.Errorf("Unexpected nil for window")
			return nil
		}

		// Check that the controls that were mounted match with the list
		if len(window.children()) != len(widgets) {
			t.Errorf("Want len(window.Children())!=nil")
		}

		err = window.SetChild(&VBox{Children: update})
		if err != nil {
			t.Errorf("Failed to set children, %s", err)
			return nil
		}

		// Check that the controls that were mounted match with the list
		if children := window.children(); children != nil {
			if len(children) != len(update) {
				t.Errorf("Wanted len(children) == len(widgets), got %d and %d", len(children), len(widgets))
			} else {
				for i := range children {
					if n1, n2 := children[i].Kind(), update[i].Kind(); n1 != n2 {
						t.Errorf("Wanted children[%d].Kind() == update[%d].Kind(), got %s and %s", i, i, n1, n2)
					} else if widget, ok := children[i].(Proper); ok {
						data := widget.Props()
						if n1, n2 := data.Kind(), update[i].Kind(); n1 != n2 {
							t.Errorf("Wanted data.Kind() == update[%d].Kind(), got %s and %s", i, n1, n2)
						}
						if !equal(t, data, update[i]) {
							t.Errorf("Wanted data == update[%d], got %v and %v", i, data, update[i])
						}
					} else {
						t.Logf("Cannot verify props of child")
					}
				}
			}
		} else {
			t.Errorf("Want window.Children()!=nil")
		}

		go func(window *Window) {
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
}

func testingUpdateWidget(t *testing.T) (updater func(base.Widget) bool, closer func()) {
	ready := make(chan *Window, 1)
	done := make(chan struct{})

	go func() {
		init := func() error {
			// Create the window.  Some of the tests here are not expected in
			// production code, but we can be a little paranoid here.
			window, err := NewWindow(t.Name(), nil)
			if err != nil {
				t.Errorf("Failed to create window, %s", err)
				return nil
			}
			if window == nil {
				t.Errorf("Unexpected nil for window")
				return nil
			}

			// Check that the controls that were mounted match with the list
			if len(window.children()) != 0 {
				t.Errorf("Want len(window.Children())!=0")
			}

			ready <- window
			return nil
		}

		err := loop.Run(init)
		if err != nil {
			t.Errorf("Failed to run GUI loop, %s", err)
		}
		close(done)
	}()

	window := <-ready

	updater = func(w base.Widget) bool {
		err := loop.Do(func() error {
			err := window.SetChild(w)
			if err != nil {
				return err
			}

			child := window.Child()
			if n1, n2 := child.Kind(), w.Kind(); n1 != n2 {
				return errors.New("child's kind does not match widget's kind")
			} else if widget, ok := child.(Proper); ok {
				data := widget.Props()
				if n1, n2 := data.Kind(), w.Kind(); n1 != n2 {
					return errors.New("child's prop's kind does not match widget's kind")
				}
				if !equal(t, data, w) {
					return errors.New("child's prop's not equal to widget")
				}
			} else {
				return errors.New("child does not support props")
			}
			return nil
		})

		if err != nil {
			t.Errorf("Error during widget update, %s", err)
			return false
		}
		return true
	}
	closer = func() {
		// Close the window
		err := loop.Do(func() error {
			window.Close()
			return nil
		})
		if err != nil {
			t.Errorf("Error in Do, %s", err)
		}

		// Wait for the GUI loop to terminate
		<-done
	}
	return updater, closer
}
