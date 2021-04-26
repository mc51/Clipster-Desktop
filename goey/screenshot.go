package goey

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"clipster/goey/loop"
)

type screenshoter interface {
	Screenshot() (image.Image, error)
}

func saveScreenshot(filename string, ss screenshoter) error {
	img, err := ss.Screenshot()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() {
		// Ignoring error when closing a file.  Possible loss of information.
		// The screenshot does not contain user data, and exists just for
		// testing.  Risk of data corruption is accepted here.
		_ = file.Close()
	}()

	return png.Encode(file, img)
}

func asyncScreenshot(filename string, window *Window) {
	// This weirdness below is to check at runtime if the type *Window
	// supports Screenshot, whose existence is platform dependent.
	ss, ok := interface{}(window).(screenshoter)
	if !ok {
		fmt.Println("error: GOEY_SCREENSHOT: Screenshots not supported on this platform.")
		return
	}

	go func() {
		// Provide a delay for the window to finish any animations.
		time.Sleep(500 * time.Millisecond)

		err := loop.Do(func() error {
			return saveScreenshot(filename, ss)
		})
		if err != nil {
			fmt.Println("error: GOEY_SCREENSHOT:", err.Error())
		}

		err = loop.Do(func() error {
			window.Close()
			return nil
		})
		if err != nil {
			fmt.Println("error: GOEY_SCREENSHOT:", err.Error())
		}
	}()
}
