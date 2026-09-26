// Functions dealing with the GUI over gotk3
package clipster

import (
	"errors"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/gen2brain/beeep"
)

var (
	server      string
	username    string
	password    string
	ssl_disable bool
)

// ShowNotification shows a desktop notification for all platforms
func ShowNotification(title string, body string) {
	// TODO: Icon in MacOS is default -> I guess it display bundle icon when there is one
	body = truncate(body, MAX_NOTIFICATION_LENGTH)
	if err := beeep.Notify(title, body, ICON_PNG_BYTES); err != nil {
		log.Println(err)
	}
}

// GUI_ConfigWindow displays the window for editing the configuration
func GUI_ConfigWindow() {
	builder, err := gtk.BuilderNewFromString(GLADE_LAYOUT)
	errorCheck(err)

	obj, err := builder.GetObject("win_creds")
	errorCheck(err)
	w, err := isWindow(obj)
	errorCheck(err)

	// Map the handlers to callback functions, and connect the signals to the Builder
	signals := map[string]interface{}{
		"form_server_address_changed_cb": onServerChange,
		"form_disable_ssl_toggled_cb":    onSSLToggle,
		"form_username_changed_cb":       onUsernameChange,
		"form_password_changed_cb":       onPasswordChange,
		"btn_login_cred_clicked_cb":      func() { onLoginBtn(w) },
		"btn_register_cred_clicked_cb":   func() { onRegisterBtn(w) },
		"btn_cancel_cred_clicked_cb":     func() { w.Close() },
	}
	builder.ConnectSignals(signals)

	// Pre-fill form with current config. Setting values triggers the signal handlers above
	server, username, password, ssl_disable = "", "", "", false
	if obj, err := builder.GetObject("form_server_address"); err == nil {
		if entry, ok := obj.(*gtk.Entry); ok && conf.Server != HOST_DEFAULT {
			entry.SetText(conf.Server)
		}
	}
	if obj, err := builder.GetObject("form_username"); err == nil {
		if entry, ok := obj.(*gtk.Entry); ok {
			entry.SetText(conf.Username)
		}
	}
	if obj, err := builder.GetObject("form_disable_ssl"); err == nil {
		if check, ok := obj.(*gtk.CheckButton); ok {
			check.SetActive(conf.Disable_ssl_cert_check)
		}
	}

	w.SetTitle("Clipster - Config")
	w.SetIcon(ICON_PNG_PIXBUF)
	w.ShowAll()
}

// GUI_FileChooserDialog displays the dialog for saving a file to disk
// and returns chosen filepath
func GUI_FileChooserDialog() string {
	var filename string
	title := "Clipster - Chose save location"
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(title, nil, gtk.FILE_CHOOSER_ACTION_SAVE,
		"_Cancel", gtk.RESPONSE_CANCEL,
		"_Save", gtk.RESPONSE_ACCEPT)
	errorCheck(err)

	dialog.SetIcon(ICON_PNG_PIXBUF)
	dialog.SetDoOverwriteConfirmation(true)
	dialog.SetCurrentName(DEFAULT_IMAGE_SAVE_NAME)
	response := dialog.Run()

	if response == gtk.RESPONSE_ACCEPT {
		filename = dialog.GetFilename()
		log.Println("Filename", filename)
	}
	dialog.Destroy()
	return filename
}

// GUI_AllClips displays the window containing all retrieved clips
func GUI_AllClips(clips []Clips) {
	builder, err := gtk.BuilderNewFromString(GLADE_LAYOUT)
	errorCheck(err)

	obj, err := builder.GetObject("win_clips")
	errorCheck(err)
	w, err := isWindow(obj)
	errorCheck(err)

	obj, err = builder.GetObject("list_clips")
	errorCheck(err)
	box, err := isListBox(obj)
	errorCheck(err)

	// Add clips to GUI Rows
	for _, clip := range clips {
		row, _ := gtk.ListBoxRowNew()
		row.SetSizeRequest(-1, 100)
		// Create rows with content
		if clip.Format == "img" {
			img, err := thumbnailToGtkImage(clip.ThumbBytes)
			if err != nil {
				log.Println("Error: creating thumbnail", err)
				continue
			}
			row.Add(img)
		} else {
			txt, _ := gtk.TextViewNew()
			txt.SetEditable(false)
			txt.SetWrapMode(gtk.WRAP_WORD_CHAR)
			buffer, _ := txt.GetBuffer()
			buffer.SetText(clip.TextDecrypted)
			row.Add(txt)
		}
		box.Add(row)
	}

	// returns the clip of the currently selected row, or false if none is selected
	selectedClip := func() (Clips, bool) {
		row := box.GetSelectedRow()
		if row == nil || row.GetIndex() < 0 || row.GetIndex() >= len(clips) {
			GUI_DialogError("Please select a clip first")
			return Clips{}, false
		}
		return clips[row.GetIndex()], true
	}

	// Map the handlers to callback functions, and connect the signals to the Builder
	signals := map[string]interface{}{
		"btn_copy_clicked_cb": func() {
			if clip, ok := selectedClip(); ok {
				SetClipboard(clip)
			}
		},
		"btn_save_clicked_cb": func() {
			if clip, ok := selectedClip(); ok {
				ImageToDisk(clip)
			}
		},
		"btn_cancel_clicked_cb": func() { w.Close() },
	}
	builder.ConnectSignals(signals)

	w.SetIcon(ICON_PNG_PIXBUF)
	w.SetTitle("Clipster - Your Clips")
	w.ShowAll()
}

// thumbnailToGtkImage creates a gtk.Image from PNG bytes. Must run on the GTK main thread
func thumbnailToGtkImage(thumb []byte) (*gtk.Image, error) {
	if len(thumb) == 0 {
		thumb = PNG_BYTES_IMAGE_NOTFOUND
	}
	pixbuf, err := gdk.PixbufNewFromBytesOnly(thumb)
	if err != nil {
		return nil, err
	}
	return gtk.ImageNewFromPixbuf(pixbuf)
}

// showMessageDialog displays a modal message dialog and blocks until it is closed
func showMessageDialog(title string, mType gtk.MessageType, buttons gtk.ButtonsType, body string) {
	// hidden parent window so GTK does not complain about a dialog without transient parent
	w, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Println("Unable to create window:", err)
		return
	}
	defer w.Destroy()
	w.SetTitle(title)
	w.SetIcon(ICON_PNG_PIXBUF)

	msg := gtk.MessageDialogNew(w, gtk.DIALOG_DESTROY_WITH_PARENT, mType, buttons, "%s", body)
	msg.SetTitle(title)
	msg.Run()
	msg.Destroy()
}

// GUI_DialogError displays an error message dialog
func GUI_DialogError(body string) {
	showMessageDialog("Clipster - Error", gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, body)
}

// GUI_DialogInfo displays an info message dialog
func GUI_DialogInfo(body string) {
	showMessageDialog("Clipster - Info", gtk.MESSAGE_INFO, gtk.BUTTONS_OK, body)
}

func onServerChange(txt *gtk.Entry) {
	var err error
	server, err = txt.GetText()
	if err != nil {
		log.Println("onServerChange", err)
	}
}

func onSSLToggle(check *gtk.CheckButton) {
	ssl_disable = check.GetActive()
}

func onUsernameChange(txt *gtk.Entry) {
	var err error
	username, err = txt.GetText()
	if err != nil {
		log.Println("onUsernameChange", err)
	}
}

func onPasswordChange(txt *gtk.Entry) {
	var err error
	password, err = txt.GetText()
	if err != nil {
		log.Println("onPasswordChange", err)
	}
}

func onLoginBtn(w *gtk.Window) {
	if err := login_flow(server, username, password, ssl_disable); err == nil {
		w.Close()
	}
}

func onRegisterBtn(w *gtk.Window) {
	if err := register_flow(server, username, password, ssl_disable); err == nil {
		w.Close()
	}
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func isListBox(obj glib.IObject) (*gtk.ListBox, error) {
	if box, ok := obj.(*gtk.ListBox); ok {
		return box, nil
	}
	return nil, errors.New("not a *gtk.ListBox")
}

func errorCheck(e error) {
	if e != nil {
		log.Panic("Gotk3 error:", e)
	}
}
