// Functions dealing with the GUI over Fyne
package clipster

import (
	"errors"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"clipster/assets"
)

const (
	textRowMinHeight   float32 = 100
	textRowMaxHeight   float32 = 300
	imageRowHeight     float32 = THUMBNAIL_HEIGHT + 20
	textRowCharsPerRow         = 80
	textRowLineHeight  float32 = 24
)

var (
	// fyneApp is set by Init. All windows are created from it
	fyneApp fyne.App
	// configWin is the currently open config window or nil. Only touch it on the main goroutine
	configWin fyne.Window
)

// sanitizeNotification makes text safe to hand to Fyne's notifications.
// On Windows Fyne pastes the text into a double quoted PowerShell string. Without this,
// a clip containing e.g. $(...), a backtick or a quote could run code. Neutralize all of them
func sanitizeNotification(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r == '"' || r == '`' || (r >= 0x2018 && r <= 0x201F): // all quote characters of PowerShell
			return '\''
		case r == '$':
			return '＄'
		case r == '\\':
			return '/'
		case !strconv.IsPrint(r): // newlines, control characters. Go would escape them
			return ' '
		}
		return r
	}, s)
}

// ShowNotification shows a desktop notification for all platforms
func ShowNotification(title string, body string) {
	body = truncate(body, MAX_NOTIFICATION_LENGTH)
	if runtime.GOOS == "windows" {
		title, body = sanitizeNotification(title), sanitizeNotification(body)
	}
	if fyneApp == nil {
		log.Println("Notification:", title, body)
		return
	}
	// safe to call from any goroutine
	fyneApp.SendNotification(fyne.NewNotification(title, body))
}

// showError displays an error message dialog over window w
func showError(w fyne.Window, err error) {
	dialog.ShowError(err, w)
}

// GUI_ConfigWindow displays the window for editing the configuration.
// Must be called on the main goroutine
func GUI_ConfigWindow() {
	if configWin != nil { // do not open it twice
		configWin.Show()
		configWin.RequestFocus()
		return
	}
	w := fyneApp.NewWindow("Clipster - Config")
	configWin = w
	w.SetOnClosed(func() { configWin = nil })
	w.SetIcon(assets.Icon)

	// Pre-fill form with current config
	server := widget.NewEntry()
	server.SetPlaceHolder(HOST_DEFAULT)
	if conf.Server != HOST_DEFAULT {
		server.SetText(conf.Server)
	}
	user := widget.NewEntry()
	user.SetText(conf.Username)
	password := widget.NewPasswordEntry()
	noSSLCheck := widget.NewCheck("", nil)
	noSSLCheck.SetChecked(conf.Disable_ssl_cert_check)

	form := widget.NewForm(
		widget.NewFormItem("Server address:", server),
		widget.NewFormItem("No SSL certification check:", noSSLCheck),
		widget.NewFormItem("Username:", user),
		widget.NewFormItem("Password:", password),
	)

	var loginBtn, registerBtn *widget.Button
	// run executes flow in the background so that the GUI does not block on the network
	run := func(flow func(string, string, string, bool) (string, error)) {
		host, name, pw, noSSL := server.Text, user.Text, password.Text, noSSLCheck.Checked
		loginBtn.Disable()
		registerBtn.Disable()
		go func() {
			msg, err := flow(host, name, pw, noSSL)
			fyne.Do(func() {
				loginBtn.Enable()
				registerBtn.Enable()
				if err != nil {
					showError(w, err)
					return
				}
				d := dialog.NewInformation("Clipster - Info", msg, w)
				d.SetOnClosed(w.Close)
				d.Show()
			})
		}()
	}
	loginBtn = widget.NewButton("Login", func() { run(login_flow) })
	registerBtn = widget.NewButton("Register", func() { run(register_flow) })
	cancelBtn := widget.NewButton("Cancel", w.Close)

	w.SetContent(container.NewPadded(container.NewVBox(
		form,
		container.NewHBox(layout.NewSpacer(), cancelBtn, registerBtn, loginBtn),
	)))
	w.Resize(fyne.NewSize(480, w.Content().MinSize().Height))
	w.CenterOnScreen()
	w.Show()
}

// GUI_AllClips displays the window containing all retrieved clips.
// Must be called on the main goroutine
func GUI_AllClips(clips []Clips) {
	w := fyneApp.NewWindow("Clipster - Your Clips")
	w.SetIcon(assets.Icon)

	// Prepare the images up front, so that updating a row stays cheap
	thumbs := make([]fyne.Resource, len(clips))
	for i, clip := range clips {
		if clip.Format != "img" {
			continue
		}
		thumb := clip.ThumbBytes
		if len(thumb) == 0 {
			thumb = assets.NotFoundPNG
		}
		thumbs[i] = fyne.NewStaticResource("thumb.png", thumb)
	}

	selected := widget.ListItemID(-1)
	list := widget.NewList(
		func() int { return len(clips) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrapWord
			img := canvas.NewImageFromResource(nil)
			img.FillMode = canvas.ImageFillContain
			return container.NewStack(label, img)
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			row := o.(*fyne.Container)
			label, img := row.Objects[0].(*widget.Label), row.Objects[1].(*canvas.Image)
			if thumbs[id] != nil {
				label.Hide()
				img.Resource = thumbs[id]
				img.Show()
				img.Refresh()
				return
			}
			img.Hide()
			label.SetText(clips[id].TextDecrypted)
			label.Show()
		},
	)
	for i, clip := range clips {
		if thumbs[i] != nil {
			list.SetItemHeight(i, imageRowHeight)
		} else {
			list.SetItemHeight(i, textRowHeight(clip.TextDecrypted))
		}
	}
	list.OnSelected = func(id widget.ListItemID) { selected = id }
	list.OnUnselected = func(widget.ListItemID) { selected = -1 }

	// returns the clip of the currently selected row, or false if none is selected
	selectedClip := func() (Clips, bool) {
		if selected < 0 || selected >= len(clips) {
			showError(w, errors.New("Please select a clip first"))
			return Clips{}, false
		}
		return clips[selected], true
	}

	copyBtn := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
		if clip, ok := selectedClip(); ok {
			go SetClipboard(clip)
		}
	})
	saveBtn := widget.NewButtonWithIcon("Save as", theme.DocumentSaveIcon(), func() {
		if clip, ok := selectedClip(); ok {
			saveImage(w, clip)
		}
	})
	closeBtn := widget.NewButton("Cancel", w.Close)

	w.SetContent(container.NewBorder(nil,
		container.NewPadded(container.NewGridWithColumns(3, copyBtn, saveBtn, closeBtn)),
		nil, nil, list))
	w.Resize(fyne.NewSize(700, 600))
	w.CenterOnScreen()
	w.Show()
}

// textRowHeight estimates how much room the text of a clip needs when displayed wrapped
func textRowHeight(text string) float32 {
	lines := 0
	for _, line := range strings.Split(text, "\n") {
		lines += max(1, int(math.Ceil(float64(len([]rune(line)))/textRowCharsPerRow)))
	}
	h := float32(lines)*textRowLineHeight + 16
	return min(max(h, textRowMinHeight), textRowMaxHeight)
}

// saveImage shows a file dialog over w and saves the image of clip at the chosen path
func saveImage(w fyne.Window, clip Clips) {
	if clip.Format != "img" {
		showError(w, errors.New("You can only save images to file!"))
		return
	}
	d := dialog.NewFileSave(func(wc fyne.URIWriteCloser, err error) {
		if err != nil {
			log.Println("Error: saving file", err)
			showError(w, err)
			return
		}
		if wc == nil { // dialog was cancelled
			return
		}
		defer wc.Close()
		path := wc.URI().Path()
		if _, err := wc.Write(clip.ImageBytes); err != nil {
			log.Println("Error: saving file", err)
			showError(w, errors.New("Error saving file: "+path+"\n"+err.Error()))
			return
		}
		log.Println("Saved file: " + path)
		dialog.ShowInformation("Clipster - Info", "File saved: "+path, w)
	}, w)
	d.SetFileName(DEFAULT_IMAGE_SAVE_NAME)
	d.Resize(w.Canvas().Size())
	d.Show()
}
