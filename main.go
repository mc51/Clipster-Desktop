// This is the main package for the Clipster-Desktop utility GUI
package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"clipster/assets"
	"clipster/internal/clipster"
)

const appID = "cc.clipster.desktop"

func main() {
	a := app.NewWithID(appID)
	a.SetIcon(assets.Icon)
	clipster.Init(a)

	desk, ok := a.(desktop.App)
	if !ok {
		log.Fatal("Error: system tray is not supported on this platform")
	}
	desk.SetSystemTrayIcon(assets.Tray)
	desk.SetSystemTrayMenu(trayMenu(desk))

	// Runs on the main goroutine once the event loop is up
	a.Lifecycle().SetOnStarted(func() {
		if _, err := clipster.LoadConfigFromFile(); err != nil {
			log.Println("Error:", err)
			clipster.GUI_ConfigWindow()
		}
	})

	a.Run()
}

// trayMenu returns the tray menu. Menu actions run on the main goroutine, so everything
// that takes time is started in a goroutine to keep the GUI responsive
func trayMenu(desk desktop.App) *fyne.Menu {
	return fyne.NewMenu("Clipster",
		fyne.NewMenuItem("Get last Clip", func() {
			log.Println("Get last Clip")
			go clipster.DownloadClipsFlow(true)
		}),
		fyne.NewMenuItem("Get all Clips", func() {
			log.Println("Get all Clips")
			go clipster.DownloadClipsFlow(false)
		}),
		fyne.NewMenuItem("Share Clip", func() {
			log.Println("Share Clip")
			go clipster.ShareClipFlow()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Edit Credentials", func() {
			log.Println("Edit Creds")
			clipster.GUI_ConfigWindow()
		}),
		// Fyne draws a check mark in front of the item as long as Checked is set
		&fyne.MenuItem{
			Label:   "Autostart Clipster",
			Checked: clipster.IsAutostartEnabled(),
			Action: func() {
				log.Println("Autostart")
				go func() {
					clipster.ToggleAutostart()
					// the tray menu can only be updated by setting a new one
					fyne.Do(func() { desk.SetSystemTrayMenu(trayMenu(desk)) })
				}()
			},
		},
	)
}
