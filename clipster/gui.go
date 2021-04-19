// Show GUI
package clipster

import (
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
	// Hidden Window Keeps GTK main loop running
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

func ShowEditCredsGUI() {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		loop.Do(GUIAskForCredentials)
	} else if runtime.GOOS == "windows" {
		StartGUIInForeground()
	}
}

// func updateWindow() {
// 	err := mainWindow.SetChild(renderCredsWindow())
// 	if err != nil {
// 		log.Println("Error:", err)
// 	}
// }

func register_flow(host string, user string, pw string) {
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
	log.Println("Registration:", host, user, pw)

	hash_login := GetLoginHashFromPw(user, pw)
	// TODO: This is blocking. Goroutine?
	if err := APIRegister(host, user, hash_login); err != nil {
		log.Println("Error:", err)
		mainWindow.Message(err.Error()).WithError().Show()
		return
	}
	hash_msg := GetMsgHashFromPw(user, pw)
	// TODO: get checkbox value
	c := Config{host, user, hash_login, hash_msg, true}
	WriteConfigFile(c)
	log.Println("Ok: Registration flow completed")
	mainWindow.Message("Registration successfull\nCredentials saved to config:\n" + CONFIG_FILEPATH).WithInfo().Show()
	mainWindow.Close()
}

func login_flow(host string, user string, pw string) {
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
	log.Println("Login:", host, user, pw)

	hash_login := GetLoginHashFromPw(user, pw)
	// TODO: This is blocking. Goroutine?
	if err := APILogin(host, user, hash_login); err != nil {
		log.Println("Error:", err)
		mainWindow.Message(err.Error()).WithError().Show()
		return
	}
	hash_msg := GetMsgHashFromPw(user, pw)
	// TODO: get checkbox value
	c := Config{host, user, hash_login, hash_msg, true}
	WriteConfigFile(c)
	log.Println("Ok: login workflow completed")
	mainWindow.Message("Login successfull\nCredentials saved to config:\n" + CONFIG_FILEPATH).WithInfo().Show()
	mainWindow.Close()
}

func renderCredsWindow() base.Widget {
	var user, pw, host string
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
						&goey.Checkbox{Text: "Disable SSL cert check",
							OnChange: func(v bool) { log.Println("check box input: ", v) }},
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
									login_flow(host, user, pw)
								}},
								&goey.Button{Text: "Register", OnClick: func() {
									register_flow(host, user, pw)
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
