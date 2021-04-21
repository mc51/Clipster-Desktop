// Show GUI
package clipster

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"guitest/goey"
	"guitest/goey/base"
	"guitest/goey/loop"

	"github.com/gen2brain/beeep"
)

var (
	mainWindow *goey.Window
)

func ShowNotification(title string, body string) {
	err := beeep.Notify(title, body, ICON_FILENAME)
	if err != nil {
		log.Println(err)
	}
}

func StartGUIInBackground() {
	// For GTK
	log.Println("StartGui")
	os.Setenv("GOEY_SIZE", "300x300") // MainWindows size
	err := loop.Run(createHiddenWindow)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println("End StartGUIInBackground")
}

func StartGUIInForeground() {
	// For Windows
	log.Println("StartGui")
	os.Setenv("GOEY_SIZE", "300x300") // MainWindows size
	err := loop.Run(GUIAskForCredentials)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println("End StartGUIInForeground")
}

func createHiddenWindow() error {
	// createHiddenWindow keeps GTK main loop running without showing a window
	// to keep tray icon available
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
	w, err := goey.NewWindow("Clipster – Enter Credentials", renderCredsWindow())
	if err != nil {
		return err
	}
	icon, err := ReadIconAsImageFromFile(ICON_FILENAME)
	if err != nil {
		return err
	}
	w.SetScroll(false, false)
	w.SetIcon(icon)
	mainWindow = w
	return nil
}

func GUIShowClips(clips []string) error {
	log.Println("GUIShowClips")
	w, err := goey.NewWindow("Clipster – Your Clips", renderShowClipsWindows(clips))
	if err != nil {
		return err
	}
	icon, err := ReadIconAsImageFromFile(ICON_FILENAME)
	if err != nil {
		return err
	}
	w.SetScroll(false, false)
	w.SetIcon(icon)
	mainWindow = w
	return nil
}

func ShowEditCredsGUI() {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		// need to run on main Thread (=GUI Thread)
		loop.Do(GUIAskForCredentials)
	} else if runtime.GOOS == "windows" {
		// can run separately
		StartGUIInForeground()
	}
}

// func updateWindow() {
// 	err := mainWindow.SetChild(renderCredsWindow())
// 	if err != nil {
// 		log.Println("Error:", err)
// 	}
// }

func DownloadAllClipsFlow() {
	// DownloadAllClipsFlow downloads all clips as json from API
	// unencrypts the encrypted texts
	// display result in gui
	clips_ecrypted, err := APIDownloadAllClips()
	if err != nil {
		mainWindow.Message(err.Error()).WithError().Show()
		log.Println("Error:", err)
		return
	}
	log.Printf("Clips: %+v", clips_ecrypted)

	clips_decrypted := make([]string, len(clips_ecrypted))

	for i := range clips_ecrypted {
		clips_decrypted[i] = Decrypt(clips_ecrypted[i].Text, HASH_ITERS_MSG)
	}
	log.Println("Unencrypted clip:", clips_decrypted)

	f := func() error { return GUIShowClips(clips_decrypted) }
	loop.Do(f)

}

func register_flow(host string, user string, pw string, ssl_disable bool) {
	// register_flow check for completeness of creds
	// creates hash from them
	// uses hash to register at API endpoint
	// displays Message box with the result
	// on success saves credentials to config
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
	// TODO: get checkbox value
	conf := Config{host, user, hash_login, hash_msg, ssl_disable}
	WriteConfigFile(conf)
	log.Println("Ok: Registration flow completed")
	mainWindow.Message("Registration successfull\nCredentials saved to config:\n" + CONFIG_FILEPATH).WithInfo().Show()
	mainWindow.Close()
}

func login_flow(host string, user string, pw string, ssl_disable bool) {
	// login_flow check for completeness of creds
	// creates hash from them
	// uses hash to authemtocate against API endpoint
	// displays Message box with the result
	// on success saves credentials to config
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
	// TODO: get checkbox value
	conf := Config{host, user, hash_login, hash_msg, ssl_disable}
	WriteConfigFile(conf)
	log.Println("Ok: login workflow completed")
	mainWindow.Message("Login successfull\nCredentials saved to config:\n" + CONFIG_FILEPATH).WithInfo().Show()
	mainWindow.Close()
}

func renderCredsWindow() base.Widget {
	// renderCredsWindow renders the Window for editing Credentials
	var user, pw, host string
	var ssl_disable bool
	widget :=
		&goey.HBox{
			Children: []base.Widget{
				&goey.VBox{
					Children: []base.Widget{
						&goey.Label{Text: "Server address:"},
						&goey.TextInput{Value: host, Placeholder: HOST_DEFAULT,
							OnChange: func(v string) {
								host = v
								log.Println("server input ", v)
							}},
						&goey.Checkbox{Text: "Disable SSL cert check", Value: false,
							OnChange: func(val_check bool) {
								ssl_disable = val_check
								log.Println("check box input: ", val_check)
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
						&goey.HBox{
							Children: []base.Widget{
								&goey.Button{Text: "Login", OnClick: func() {
									login_flow(host, user, pw, ssl_disable)
								}},
								&goey.Button{Text: "Register", OnClick: func() {
									register_flow(host, user, pw, ssl_disable)
								}},
								&goey.Button{Text: "Cancel", OnClick: func() { mainWindow.Close() }}},
							AlignMain: goey.MainStart,
						},
					},
				},
			},
			AlignMain: goey.MainCenter,
		}
	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child:  widget,
	}
}

func renderShowClipsWindows(clips []string) base.Widget {
	// renderShowClipsWindows renders goey Window showing downloaded Clips
	widgets := []base.Widget{
		&goey.Label{Text: "Here are your Clips:"},
	}

	for _, v := range clips {
		widgets = append(widgets, &goey.TextInput{Value: v,
			OnChange:   func(v string) { println("text input ", v) },
			OnEnterKey: func(v string) { println("t1* ", v) }})
	}

	widgets = append(widgets, &goey.HBox{
		Children: []base.Widget{
			&goey.Button{Text: "Copy to Clipboatd", OnClick: func() {
				fmt.Println("Copy to Clipboard")
			}},
			&goey.Button{Text: "Cancel", OnClick: func() { mainWindow.Close() }}},
		AlignMain: goey.MainStart,
	})

	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child:  &goey.VBox{Children: widgets},
	}
}
