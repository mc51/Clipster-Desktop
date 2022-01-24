// This is the main package for the Clipster-Desktop utility GUI
package main

import (
	"log"
	"os"

	"clipster/clipster"

	"github.com/faiface/mainthread"
	"github.com/getlantern/systray"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	mainthread.Run(run) // enables mainthread package and runs run in a separate goroutine
}

// run GUI on mainthread. If using GTK wait for main loop start
// then check for Config
func run() {
	finish := make(chan bool)
	// On MacOS GUI needs to be running on main thread or we get a panic
	systray.Register(onReady, onExit)
	mainthread.CallNonBlock(func() { startGui(finish) })
	// For GTK wait until main loop is started
	go func() {
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

// startGui starts systray and GUI loops. It deals with platform idiosyncrasies
func startGui(finish chan bool) {
	gtk.Init(nil)
	gtk.Main()
	close(finish)
}

// onReady is called on systray startup. It displays tray menu and deals with selections
func onReady() {
	log.Println("On Ready")
	systray.SetIcon(clipster.ICON_TRAY_BYTES)
	systray.SetTitle("Clipster")
	systray.SetTooltip("Clipster")
	autostart_enabled := clipster.IsAutostartEnabled()

	// We can manipulate the tray in other goroutines
	go func() {

		mLastClip := systray.AddMenuItem("Get last Clip", "Get last Clip")
		mAllClips := systray.AddMenuItem("Get all Clips", "Get all Clips")
		mShareClip := systray.AddMenuItem("Share Clip", "Share Clip")
		systray.AddSeparator()
		mEditCreds := systray.AddMenuItem("Edit Credentials", "Edit Credentials")
		mAutostart := systray.AddMenuItemCheckbox("Autostart Clipster", "Autostart Clipster",
			autostart_enabled)
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

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
				// TODO: FIXME this doesnt work on windows - checkmark status changed only
				// after restart
				autostart_enabled = !autostart_enabled
				clipster.ToggleAutostart()
			case <-mQuit.ClickedCh:
				log.Println("Quit")
				onExit()
				return
			}
		}
	}()
}

// onExit is called on systray menu quit
func onExit() {
	// Remove temp icon file
	if err := os.Remove(clipster.ICON_FILENAME); err != nil {
		log.Println("Error: deleting temp file", err)
	}
	systray.Quit()
}
