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
			time.Sleep(200 * time.Millisecond)
			log.Println("Waiting for GTK loop to start...")
		}
		log.Printf("GTK loop started... Checking for config")
		ok, err := clipster.OpenConfigFile()
		if !ok {
			log.Println("Error:", err)
			clipster.ShowEditCredsGUI()
			log.Println("After Show Gui")
		} else {
			creds, _ := clipster.LoadConfigFromFile()
			log.Printf("%+v", creds)
		}
	}()
	<-finish
}

func startGui(finish chan bool) {
	// startGui starts systray and GUI loops. It deals with platform idiosyncraties
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		// On linux and macos all GUIs must run on main thead. We use GTK for tray and goey
		// Both must be run in same loop, locked to main thread
		tray.Register(onReady, onExit)
		// Start gtk loop without displaying window (just tray)
		clipster.StartGUIInBackground()
	} else if runtime.GOOS == "windows" {
		// On Win GUIs can run on non-main thread. tray gets own thread
		tray.Run(onReady, onExit)
	}
	close(finish)
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
				clipster.DownloadAllClipsFlow()
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
