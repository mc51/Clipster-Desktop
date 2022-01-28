// Functions dealing with the GUI over gotk3
package clipster

import (
	"errors"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/gen2brain/beeep"
)

var (
	sel_list_row int
	server       string
	username     string
	password     string
	ssl_disable  bool
)

func ShowNotification(title string, body string) {
	// TODO: Icon in MacOS is default -> I guess it display bundle icon when there is one
	if len(body) >= MAX_NOTIFICATION_LENGTH {
		body = body[0:MAX_NOTIFICATION_LENGTH] + " [...]"
	}
	err := beeep.Notify(title, body, ICON_FILENAME)
	if err != nil {
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

	w.SetTitle("Clipster - Config")
	w.Connect("destroy", func() {
		w.Close()
	})
	w.ShowAll()
}

// GUI_FileChooserDialog displays the dialog for saving a file to disk
// and returns chosen filepath
func GUI_FileChooserDialog() string {

	var filename string
	title := "Clipster - Chose save location"
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(title, nil, gtk.FILE_CHOOSER_ACTION_SAVE,
		"gtk-cancel", gtk.RESPONSE_CANCEL,
		"gtk-save", gtk.RESPONSE_ACCEPT)
	errorCheck(err)

	dialog.SetDoOverwriteConfirmation(true)
	dialog.SetCurrentName(DEFAULT_IMAGE_SAVE_NAME)
	response := dialog.Run()

	if response == gtk.RESPONSE_ACCEPT {
		chooser := dialog.FileChooser
		filename = chooser.GetFilename()
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
			log.Println("Got image clip")
			clip := processClipTextToImages(clip)
			row.Add(clip.GtkThumb)
		} else {
			log.Println("Got text clip")
			txt, _ := gtk.TextViewNew()
			txt.SetEditable(false)
			txt.SetWrapMode(gtk.WRAP_WORD_CHAR)
			buffer, _ := txt.GetBuffer()
			buffer.SetText(clip.TextDecrypted)
			row.Add(txt)
		}
		box.Add(row)
	}

	// Map the handlers to callback functions, and connect the signals to the Builder
	signals := map[string]interface{}{
		"list_clips_row_selected_cb": func(obj *gtk.ListBox) { sel_list_row = obj.GetSelectedRow().GetIndex() },
		"btn_copy_clicked_cb":        func(obj *gtk.Button) { SetClipboard(clips[sel_list_row]) },
		"btn_save_clicked_cb":        func(obj *gtk.Button) { ImageToDisk(clips[sel_list_row]) },
		"btn_cancel_clicked_cb":      func() { w.Close() },
	}
	builder.ConnectSignals(signals)

	w.SetTitle("Clipster - Your Clips")
	w.Connect("destroy", func() {
		w.Close()
	})
	w.ShowAll()
}

func createWindow(title string) *gtk.Window {
	w, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatalln("Unable to create window:", err)
		return nil
	}
	w.SetTitle(title)
	w.Connect("destroy", func() {
		w.Close()
	})
	return w
}

// GUI_DialogError displays an error message dialog
func GUI_DialogError(body string) {
	w := createWindow("Clipster - Error")
	msg := gtk.MessageDialogNew(w, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR,
		gtk.BUTTONS_CLOSE, body)
	msg.Connect("response", func() { msg.Destroy() })
	msg.Run()
}

// GUI_DialogInfo displays an info message dialog
func GUI_DialogInfo(body string) {
	w := createWindow("Clipster - Info")
	msg := gtk.MessageDialogNew(w, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO,
		gtk.BUTTONS_OK, body)
	msg.Connect("response", func() { msg.Destroy() })
	msg.Run()
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
	if err := login_flow(server, username, password, ssl_disable); err != nil {
		return
	} else {
		w.Close()
	}
}

func onRegisterBtn(w *gtk.Window) {
	if err := register_flow(server, username, password, ssl_disable); err != nil {
		return
	} else {
		w.Close()
	}
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func isFileChooserDialog(obj glib.IObject) (*gtk.FileChooserDialog, error) {
	if dialog, ok := obj.(*gtk.FileChooserDialog); ok {
		return dialog, nil
	}
	return nil, errors.New("not a *gtk.FileChooserDialog")
}

func isButton(obj glib.IObject) (*gtk.Button, error) {
	if btn, ok := obj.(*gtk.Button); ok {
		return btn, nil
	}
	return nil, errors.New("not a *gtk.Button")
}

func isListBox(obj glib.IObject) (*gtk.ListBox, error) {
	// Make type assertion (as per gtk.go).
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
