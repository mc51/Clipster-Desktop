// Show GUI
package clipster

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"clipster/goey"
	"clipster/goey/base"
	"clipster/goey/loop"

	"github.com/gen2brain/beeep"
)

var (
	mainWindow *goey.Window
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

// For GTK to start main loop without showing a window (to show trayicon)
func StartGUIInBackground() {
	log.Println("StartGui")
	err := loop.Run(createHiddenWindow)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println("End StartGUIInBackground")
}

// createHiddenWindow keeps GTK main loop running without showing a window
// to keep tray icon available
func createHiddenWindow() error {
	log.Println("createHiddenWindow")
	// this need adapted goey function without gtk call to .show()
	_, err := goey.NewHiddenWindow("", nil)
	if err != nil {
		log.Fatalln("Error:", err)
		return err
	}
	return nil
}

func GUIAskForCredentials() error {
	log.Println("GUIAskForCredentials")
	os.Setenv("GOEY_SIZE", "300x300")
	w, err := goey.NewWindow("Clipster – Enter Credentials", renderCredsWindow())
	if err != nil {
		return err
	}
	icon, err := IconAsImageFromBytes(ICON_PNG_BYTES)
	if err != nil {
		return err
	}
	w.SetScroll(false, false)
	w.SetIcon(icon)
	mainWindow = w
	return nil
}

func GUIShowClips(clips []Clips) error {
	log.Println("GUIShowClips")
	os.Setenv("GOEY_SIZE", "600x100")
	w, err := goey.NewWindow("Clipster – Your Clips", renderShowClipsWindows(clips))
	if err != nil {
		return err
	}
	icon, err := IconAsImageFromBytes(ICON_PNG_BYTES)
	if err != nil {
		return err
	}
	w.SetScroll(false, false)
	w.SetIcon(icon)
	mainWindow = w
	return nil
}

// guiDo runs a GUI function on the appropriate thread depending on the os
func guiDo(f func() error) {
	log.Println("guiDo")
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		// run f on main (GUI) thread
		err := loop.Do(f)
		if err != nil {
			log.Panicln("Error:", err)
		}
	} else if runtime.GOOS == "windows" {
		// start new thread for f
		err := loop.Run(f)
		if err != nil {
			log.Panicln("Error:", err)
		}
	}
}

func ShowEditCredsGUI() {
	log.Println("ShowEditCredsGUI")
	guiDo(GUIAskForCredentials)
}

// func updateWindow() {
// 	err := mainWindow.SetChild(renderCredsWindow())
// 	if err != nil {
// 		log.Println("Error:", err)
// 	}
// }

// ShareClipFlow gets current clipboard value, encrypts it
// uploads it to server and shows notification
func ShareClipFlow() {
	log.Println("ShareClipFlow")
	clip, format := GetClipboard()
	clip_encrypted := Encrypt(clip)
	if err := APIShareClip(clip_encrypted, format); err != nil {
		ShowNotification("Clipster - Error", err.Error())
		log.Println("Error:", err)
		return
	}
	if format == "txt" {
		ShowNotification("Clipster – Shared clip", clip)
	} else if format == "img" {
		ShowNotification("Clipster – Shared clip", MSG_NOTIFY_GOT_IMAGE)
	}
}

// DownloadLastClipFlow downloads all clips as json from API
// unencrypts the latest encrypted text
// copies content to clipboard and shows notification
func DownloadLastClipFlow() {
	log.Println("DownloadLastClipFlow")
	clips_encrypted, err := APIDownloadAllClips()
	if err != nil {
		ShowNotification("Clipster - Error", err.Error())
		log.Println("Error:", err)
		return
	}
	log.Printf("Clips: %+v", clips_encrypted)
	last_clip := clips_encrypted[len(clips_encrypted)-1]
	last_clip.TextDecrypted = Decrypt(last_clip.Text)
	SetClipboard(last_clip)
}

// DownloadAllClipsFlow downloads all clips as json from API, unencrypts the encrypted texts
// and display result in gui
func DownloadAllClipsFlow() {
	clips, err := APIDownloadAllClips()
	if err != nil {
		ShowNotification("Clipster - Error", err.Error())
		log.Println("Error:", err)
		return
	}
	log.Printf("Clips: %+v", clips)
	for i := range clips {
		clips[i].TextDecrypted = Decrypt(clips[i].Text)
	}
	f := func() error { return GUIShowClips(clips) }
	guiDo(f)
}

// register_flow check for completeness of creds, creates hash from them,
// uses hash to register at API endpoint, displays Message box with the result.
// On success saves credentials to config
func register_flow(host string, user string, pw string, ssl_disable bool) {
	host, user, pw, err := AreCredsComplete(host, user, pw)
	if err != nil {
		mainWindow.Message(err.Error()).WithError().Show()
		log.Println("Error:", err)
		return
	}
	// TODO: Remove all cleartext pws from logs?
	log.Println("Registration:", host, user, pw, ssl_disable)

	hash_login := GetLoginHashFromPw(user, pw)
	// TODO: This is blocking. Goroutine?
	if err := APIRegister(host, user, hash_login, ssl_disable); err != nil {
		log.Println("Error:", err)
		mainWindow.Message(err.Error()).WithError().Show()
		return
	}
	hash_msg := GetMsgHashFromPw(user, pw)
	conf = Config{host, user, hash_login, hash_msg, ssl_disable}
	WriteConfigFile(conf)
	log.Println("Ok: Registration flow completed")
	mainWindow.Message("Registration successfull\nCredentials saved to config:\n" +
		CONFIG_FILEPATH).WithInfo().Show()
	mainWindow.Close()
}

// login_flow check for completeness of creds, creates hash from them,
// uses hash to authenticate against API endpoint, displays Message box with the result.
// On success saves credentials to config
func login_flow(host string, user string, pw string, ssl_disable bool) {
	host, user, pw, err := AreCredsComplete(host, user, pw)
	if err != nil {
		mainWindow.Message(err.Error()).WithError().Show()
		log.Println("Error:", err)
		return
	}
	// TODO: Remove all cleartext pws from logs?
	log.Println("Login:", host, user, pw, ssl_disable)

	hash_login := GetLoginHashFromPw(user, pw)
	// TODO: This is blocking. Goroutine?
	if err := APILogin(host, user, hash_login, ssl_disable); err != nil {
		log.Println("Error:", err)
		mainWindow.Message(err.Error()).WithError().Show()
		return
	}
	hash_msg := GetMsgHashFromPw(user, pw)
	conf = Config{host, user, hash_login, hash_msg, ssl_disable}
	WriteConfigFile(conf)
	log.Println("Ok: login workflow completed")
	mainWindow.Message("Login successfull\nCredentials saved to config:\n" +
		CONFIG_FILEPATH).WithInfo().Show()
	mainWindow.Close()
}

// renderCredsWindow renders the Window for editing Credentials
func renderCredsWindow() base.Widget {
	var user, pw, host string
	var ssl_disable bool
	widget :=
		[]base.Widget{
			&goey.Label{Text: "Server address:"},
			&goey.TextInput{Value: host, Placeholder: HOST_DEFAULT,
				OnChange: func(v string) {
					host = v
					log.Println("server input ", v)
				}},
			&goey.Checkbox{Text: "Disable SSL cert check", Value: false,
				OnChange: func(val_check bool) {
					ssl_disable = val_check
					log.Println("SSL checkbox: ", ssl_disable)
				}},
			&goey.Label{Text: "Username:"},
			&goey.TextInput{Value: user, Placeholder: "Enter username",
				OnChange: func(v string) {
					user = v
					log.Println("username input ", v)
				}},
			&goey.Label{Text: "Password:"},
			&goey.TextInput{Value: pw, Placeholder: "Enter password", Password: true,
				OnChange: func(v string) {
					pw = v
					log.Println("password input ", v)
				}},
		}

	widget = append(widget, &goey.HBox{
		Children: []base.Widget{
			&goey.Button{Text: "Login", OnClick: func() {
				login_flow(host, user, pw, ssl_disable)
			}},
			&goey.Button{Text: "Register", OnClick: func() {
				register_flow(host, user, pw, ssl_disable)
			}},
			&goey.Button{Text: "Cancel", OnClick: func() { mainWindow.Close() }}},
		AlignMain: goey.MainCenter,
	})

	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child:  &goey.VBox{Children: widget},
	}
}

// renderShowClipsWindows renders goey Window showing downloaded Clips.
// Allows user to copy text using shortcut or by clicking into text field
// and using the button
func renderShowClipsWindows(clips []Clips) base.Widget {
	var id_selected int
	widgets := []base.Widget{
		&goey.Label{Text: "Your shared Clips:"},
	}

	// TODO: when img, show thumbnail or something
	for i, v := range clips {
		j := i
		widgets = append(widgets, &goey.TextArea{Value: v.TextDecrypted,
			ReadOnly: true,
			MinLines: 3,
			OnFocus: func() {
				id_selected = j
			}})
	}

	widgets = append(widgets, &goey.HBox{
		Children: []base.Widget{
			&goey.Button{Text: "Copy to Clipboard", OnClick: func() {
				fmt.Println("Copy to Clipboard")
				clip := clips[id_selected]
				SetClipboard(clip)
				mainWindow.Close()
			}},
			&goey.Button{Text: "Cancel", OnClick: func() { mainWindow.Close() }}},
		AlignMain: goey.MainStart,
	})

	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child:  &goey.VBox{Children: widgets},
	}
}
