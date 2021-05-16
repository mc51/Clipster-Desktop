// This is the main package for the Clipster-Desktop utility GUI
package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"clipster/clipster"
	"clipster/goey/loop"
	"clipster/tray"

	"github.com/faiface/mainthread"
)

func main() {
	mainthread.Run(run) // enables mainthread package and runs run in a separate goroutine
}

func run() {
	finish := make(chan bool)
	// On MacOS GUI needs to be running on main thread or we get a panic
	mainthread.CallNonBlock(func() { startGui(finish) })
	// For GTK wait until main loop is started
	go func() {
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			for loop.IsRunning == 0 {
				log.Println("Waiting for GTK loop to start...")
				time.Sleep(200 * time.Millisecond)
			}
			log.Printf("GTK loop started... moving on")
		}
		ok, err := clipster.OpenConfigFile()
		if !ok {
			log.Println("Error:", err)
			clipster.ShowEditCredsGUI()
		} else {
			conf, _ := clipster.LoadConfigFromFile()
			log.Printf("%+v", conf)
		}
	}()
	<-finish
}

func startGui(finish chan bool) {
	// startGui starts systray and GUI loops. It deals with platform idiosyncrasies
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		// On linux and macos all GUIs must run on single main thead.
		// We use GTK for tray and goey. Both must run in same loop, locked to main thread
		tray.Register(onReady, onExit)
		// Start gtk loop without displaying window (to show tray)
		clipster.StartGUIInBackground()
	} else if runtime.GOOS == "windows" {
		// On Win GUIs can run on non-main thread. tray gets own thread
		tray.Run(onReady, onExit)
	}
	close(finish)
}

func onReady() {
	log.Println("On Ready")
	tray.SetIcon(clipster.ICON_TRAY_BYTES)
	tray.SetTitle("Clipster")
	tray.SetTooltip("Clipster")
	autostart_enabled := clipster.IsAutostartEnabled()

	// We can manipulate the tray in other goroutines
	go func() {

		mLastClip := tray.AddMenuItem("Get last Clip", "Get last Clip")
		mAllClips := tray.AddMenuItem("Get all Clips", "Get all Clips")
		mShareClip := tray.AddMenuItem("Share Clip", "Share Clip")
		tray.AddSeparator()
		mEditCreds := tray.AddMenuItem("Edit Credentials", "Edit Credentials")
		mAutostart := tray.AddMenuItemCheckbox("Autostart Clipster", "Autostart Clipster",
			autostart_enabled)
		tray.AddSeparator()
		mQuit := tray.AddMenuItem("Quit", "Quit the whole app")

		// Read from Channel: Called as callback from C
		for {
			select {
			case <-mLastClip.ClickedCh:
				log.Println("Get last Clip")
				clipster.DownloadLastClipFlow()
			case <-mAllClips.ClickedCh:
				log.Println("Get all Clips")
				clipster.DownloadAllClipsFlow()
			case <-mShareClip.ClickedCh:
				log.Println("Share Clip")
				clipster.ShareClipFlow()
			case <-mEditCreds.ClickedCh:
				log.Println("Edit Creds")
				clipster.ShowEditCredsGUI()
			case <-mAutostart.ClickedCh:
				log.Println("Autostart")
				clipster.ToggleAutostart()
			case <-mQuit.ClickedCh:
				log.Println("Quit")
				onExit()
				return
			}
		}
	}()
}

func onExit() {
	// Remove temp icon file
	if err := os.Remove(clipster.ICON_FILENAME); err != nil {
		log.Println("Error: deleting temp file", err)
	}
	tray.Quit()
}
