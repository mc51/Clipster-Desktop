package goey

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"testing"
	"time"

	"guitest/goey/base"
	"guitest/goey/loop"
)

func ExampleNewWindow() {
	// All calls that modify GUI objects need to be schedule ont he GUI thread.
	// This callback will be used to create the top-level window.
	createWindow := func() error {
		// Create a top-level window.
		mw, err := NewWindow("Test", &VBox{
			Children: []base.Widget{
				&Button{Text: "Click me!"},
			},
		})
		if err != nil {
			// This error will be reported back up through the call to
			// Run below.  No need to print or log it here.
			return err
		}

		// We can start a goroutine, but note that we can't modify GUI objects
		// directly.
		go func() {
			fmt.Println("Up")
			time.Sleep(50 * time.Millisecond)
			fmt.Println("Down")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			_ = loop.Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	// Start the GUI thread.
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Output:
	// Up
	// Down
}

func ExampleWindow_Message() {
	// All calls that modify GUI objects need to be schedule ont he GUI thread.
	// This callback will be used to create the top-level window.
	createWindow := func() error {
		// Create a top-level window.
		mw, err := NewWindow("Test", &Button{Text: "Click me!"})
		if err != nil {
			// This error will be reported back up through the call to
			// Run below.  No need to print or log it here.
			return err
		}

		// We can start a goroutine, but note that we can't modify GUI objects
		// directly.
		go func() {
			// Show the error message.
			_ = loop.Do(func() error {
				return mw.Message("This is an example message.").WithInfo().Show()
			})

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			_ = loop.Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	// Start the GUI thread.
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func testingWindow(t *testing.T, action func(*testing.T, *Window)) {
	createWindow := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		mw, err := NewWindow(t.Name(), nil)
		if err != nil {
			t.Fatalf("failed to create window: %s", err)
		}
		if mw == nil {
			t.Fatalf("unexpected nil for window")
		}

		go func() {
			// Delegate to test specific actions.
			action(t, mw)
			if testing.Verbose() {
				time.Sleep(250 * time.Millisecond)
			}

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			err := loop.Do(func() error {
				mw.Close()
				return nil
			})
			if err != nil {
				// Would like to report this error using t.Fatalf, but we are
				// not in the same goroutine.  Could send a message using a
				// channel, but if the call to Do failed, it is not certain that
				// we closed the window, and could deadlock.
				panic(err)
			}
		}()

		return nil
	}

	err := loop.Run(createWindow)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
}

func TestMain(m *testing.M) {
	loop.TestMain(m)
}

func TestWindow_SetChild(t *testing.T) {
	testingWindow(t, func(t *testing.T, mw *Window) {
		widgets := []base.Widget{}

		for i := 1; i < 10; i++ {
			if testing.Verbose() {
				time.Sleep(250 * time.Millisecond)
			} else {
				time.Sleep(50 * time.Millisecond)
			}
			widgets = append(widgets, &Button{Text: "Button " + strconv.Itoa(i)})
			err := loop.Do(func() error {
				return mw.SetChild(&VBox{
					AlignMain:  SpaceBetween,
					AlignCross: CrossCenter,
					Children:   widgets,
				})
			})
			if err != nil {
				t.Logf("Error setting children, %s", err)
			}
		}
		for i := len(widgets); i > 0; i-- {
			time.Sleep(50 * time.Millisecond)
			widgets = widgets[:i-1]
			err := loop.Do(func() error {
				return mw.SetChild(&VBox{
					AlignMain:  SpaceBetween,
					AlignCross: CrossCenter,
					Children:   widgets,
				})
			})
			if err != nil {
				t.Logf("Error setting children, %s", err)
			}
		}
		time.Sleep(50 * time.Millisecond)
	})
}

func makeImage(t *testing.T, index int) image.Image {
	colors := [3]color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	bounds := image.Rect(0, 0, 32, 32)
	img := image.NewRGBA(bounds)
	draw.Draw(img, image.Rect(0, 0, 11, 32), image.NewUniform(colors[index%3]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(11, 0, 22, 32), image.NewUniform(colors[(index+1)%3]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(22, 0, 32, 32), image.NewUniform(colors[(index+2)%3]), image.Point{}, draw.Src)
	return img
}

func TestNewWindow_SetIcon(t *testing.T) {
	testingWindow(t, func(t *testing.T, mw *Window) {
		for i := 0; i < 6; i++ {
			img := makeImage(t, i)

			err := loop.Do(func() error {
				return mw.SetIcon(img)
			})
			if err != nil {
				t.Errorf("Error calling SetIcon, %s", err)
			}
			time.Sleep(50 * time.Millisecond)
		}
	})
}

func TestNewWindow_SetScroll(t *testing.T) {
	testingWindow(t, func(t *testing.T, mw *Window) {
		cases := []struct {
			horizontal, vertical bool
		}{
			{false, false},
			{false, true},
			{true, false},
			{true, true},
		}

		for i, v := range cases {
			err := loop.Do(func() error {
				mw.SetScroll(v.horizontal, v.vertical)
				out1, out2 := mw.Scroll()
				if out1 != v.horizontal {
					t.Errorf("Case %d: Returned flag for horizontal scroll does not match, got %v, want %v", i, out1, v.horizontal)
				}
				if out2 != v.vertical {
					t.Errorf("Case %d: Returned flag for vertical scroll does not match, got %v, want %v", i, out2, v.vertical)
				}
				return nil
			})
			if err != nil {
				t.Errorf("Error calling SetTitle, %s", err)
			}
		}
	})
}

func TestNewWindow_SetTitle(t *testing.T) {
	testingWindow(t, func(t *testing.T, mw *Window) {
		err := loop.Do(func() error {
			err := mw.SetTitle("Flash!")
			if err != nil {
				return err
			}

			if got, err := mw.Title(); err != nil {
				t.Errorf("Failed to get title, got error %s", err)
			} else if got != "Flash!" {
				t.Errorf("Failed to set title correctly, got %s", got)
			}
			return nil
		})
		if err != nil {
			t.Errorf("Error calling SetTitle, %s", err)
		}
	})
}
