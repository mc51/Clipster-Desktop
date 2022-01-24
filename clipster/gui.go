// Show GUI
package clipster

import (
	"errors"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/gen2brain/beeep"
)

var (
	selected_list_row int
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

func DownloadLastClipFlow() error {
	return nil
}

func DownloadAllClipsFlow() error {
	return nil
}

func ShareClipFlow() error {
	return nil
}

func ShowEditCredsGUI() error {
	return nil
}

func onCopyBtn(b *gtk.Button) {
	log.Println("onCopyBtn", b)
}

func onSaveBtn() {
	log.Println("onSaveBtn")
}

func onCancelBtn() {
	log.Println("onCancelBtn")
	gtk.MainQuit()
}

func onSelectedRow(listbox *gtk.ListBox) {
	log.Println("onSelectedRow", listbox)
	selected_list_row = listbox.GetSelectedRow().GetIndex()
	log.Println("onSelectedRow", selected_list_row)
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func isButton(obj glib.IObject) (*gtk.Button, error) {
	if btn, ok := obj.(*gtk.Button); ok {
		return btn, nil
	}
	return nil, errors.New("not a *gtk.Button")
}

func isBox(obj glib.IObject) (*gtk.Box, error) {
	// Make type assertion (as per gtk.go).
	if box, ok := obj.(*gtk.Box); ok {
		return box, nil
	}
	return nil, errors.New("not a *gtk.Box")
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
