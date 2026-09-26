// This is the main package for the Clipster-Desktop utility GUI
package main

import (
	_ "embed"
	"log"
	"runtime"

	"clipster/clipster"

	"github.com/getlantern/systray"
	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/clipster.glade
var GLADE_LAYOUT string

// GTK (and Cocoa on MacOS) must run on the main OS thread
func init() {
	runtime.LockOSThread()
}

func main() {
	clipster.GLADE_LAYOUT = GLADE_LAYOUT
	clipster.Init()

	gtk.Init(nil)
	systray.Register(onReady, onExit)

	if err := clipster.OpenConfigFile(); err != nil {
		log.Println("Error:", err)
		clipster.DoGUI(clipster.GUI_ConfigWindow)
	} else if _, err := clipster.LoadConfigFromFile(); err != nil {
		log.Println("Error:", err)
		clipster.DoGUI(clipster.GUI_ConfigWindow)
	}

	gtk.Main()
}

// onReady is called on systray startup. It displays tray menu and deals with selections
func onReady() {
	log.Println("On Ready")
	systray.SetIcon(clipster.ICON_TRAY_BYTES)
	systray.SetTitle("Clipster")
	systray.SetTooltip("Clipster")

	mLastClip := systray.AddMenuItem("Get last Clip", "Get last Clip")
	mAllClips := systray.AddMenuItem("Get all Clips", "Get all Clips")
	mShareClip := systray.AddMenuItem("Share Clip", "Share Clip")
	systray.AddSeparator()
	mEditCreds := systray.AddMenuItem("Edit Credentials", "Edit Credentials")
	mAutostart := systray.AddMenuItemCheckbox("Autostart Clipster", "Autostart Clipster",
		clipster.IsAutostartEnabled())
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// onReady already runs in its own goroutine, so we can block here
	for {
		select {
		case <-mLastClip.ClickedCh:
			log.Println("Get last Clip")
			clipster.DownloadClipsFlow(true)
		case <-mAllClips.ClickedCh:
			log.Println("Get all Clips")
			clipster.DownloadClipsFlow(false)
		case <-mShareClip.ClickedCh:
			log.Println("Share Clip")
			clipster.ShareClipFlow()
		case <-mEditCreds.ClickedCh:
			log.Println("Edit Creds")
			clipster.DoGUI(clipster.GUI_ConfigWindow)
		case <-mAutostart.ClickedCh:
			log.Println("Autostart")
			clipster.ToggleAutostart()
			if clipster.IsAutostartEnabled() {
				mAutostart.Check()
			} else {
				mAutostart.Uncheck()
			}
		case <-mQuit.ClickedCh:
			log.Println("Quit")
			systray.Quit()
			return
		}
	}
}

// onExit is called by systray when it shuts down
func onExit() {
	log.Println("On Exit")
}
