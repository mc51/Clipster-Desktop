// Package assets embeds the images used by Clipster
package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed Icon.png
var iconPNG []byte

//go:embed clipster_icon_128.png
var trayPNG []byte

// NotFoundPNG is a placeholder image for clips that can not be decoded
//
//go:embed not_found_64.png
var NotFoundPNG []byte

var (
	// Icon is the application icon (windows, dialogs, notifications)
	Icon fyne.Resource = fyne.NewStaticResource("Icon.png", iconPNG)
	// Tray is the monochrome icon that stays readable at tray icon sizes
	Tray fyne.Resource = fyne.NewStaticResource("tray.png", trayPNG)
)
