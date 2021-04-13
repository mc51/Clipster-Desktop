// This is the main package for the Clipster-Desktop utility GUI
package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"guitest/clipster"
	"guitest/goey/loop"
	"guitest/tray"
)

func main() {
	finish := make(chan bool)
	go startGui(finish)
	go func() {
		for loop.IsRunning == 0 {
			time.Sleep(1 * time.Second)
			log.Println("Waiting for GTK loop to start...")
		}
		log.Printf("GTK running...")
		clipster.ShowEditCredsGUI()
	}()
	<-finish
}

func startGui(finish chan bool) {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		tray.Register(onReady, onExit)
		clipster.StartGuiInBackground()
	} else if runtime.GOOS == "windows" {
		// Start systray with GTK loop and GUI
		tray.Run(onReady, onExit)
	}
	finish <- true
}

func onReady() {
	fmt.Println("On Ready")
	tray.SetIcon(clipster.ReadIconAsBytesFromFile(clipster.ICON_FILENAME))
	tray.SetTitle("Clipster")
	tray.SetTooltip("Clipster")

	// We can manipulate the tray in other goroutines
	go func() {
		mLastClip := tray.AddMenuItem("Get last Clip", "Get last Clip")
		mAllClips := tray.AddMenuItem("Get all Clips", "Get all Clips")
		mShareClip := tray.AddMenuItem("Share Clip", "Share Clip")
		tray.AddSeparator()
		mEditCreds := tray.AddMenuItem("Edit Credentials", "Edit Credentials")
		tray.AddSeparator()
		mQuit := tray.AddMenuItem("Quit", "Quit the whole app")

		// Read from Channel: Called as callback from C
		for {
			select {
			case <-mLastClip.ClickedCh:
				log.Println("Last")
			case <-mAllClips.ClickedCh:
				log.Println("All")
			case <-mShareClip.ClickedCh:
				log.Println("Share")
			case <-mEditCreds.ClickedCh:
				log.Println("Creds")
				clipster.ShowEditCredsGUI()
			case <-mQuit.ClickedCh:
				tray.Quit()
				log.Println("Exiting")
				return
			}
		}
	}()

}

func onExit() {
	tray.Quit()
}
