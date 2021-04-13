// Show GUI
package clipster

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"

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

func StartGuiInBackground() {
	// For GTK
	fmt.Println("StartGui")
	os.Setenv("GOEY_SIZE", "300x300") // MainWindows size
	err := loop.Run(createHiddenWindow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	log.Println("StartGuiInBackground")
}

func StartGuiInForeground() {
	// For Windows
	fmt.Println("StartGui")
	os.Setenv("GOEY_SIZE", "300x300") // MainWindows size
	err := loop.Run(CreateMainWindow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}

func createHiddenWindow() error {
	// Hidden Window Keeps GTK main loop running
	fmt.Println("createHiddenWindow")
	// this need adapted goey function without gtk call to .show()
	_, err := goey.NewHiddenWindow("", nil)
	if err != nil {
		return err
	}
	return nil
}

func CreateMainWindow() error {
	fmt.Println("CreateMainWindow")
	w, err := goey.NewWindow("Clipster – Enter Credentials", renderWindow())
	if err != nil {
		return err
	}
	icon, err := ReadIcon()
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
		loop.Do(CreateMainWindow)
	} else if runtime.GOOS == "windows" {
		StartGuiInForeground()
	}
}

func ReadIcon() (image.Image, error) {
	icon_b64 := ReadIconAsB64FromFile(ICON_FILENAME)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(icon_b64))
	m, _, err := image.Decode(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	return m, err
}

func updateWindow() {
	err := mainWindow.SetChild(renderWindow())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}

func login(host string, user string, pw string) {
	host, user, pw, err := areCredsComplete(host, user, pw)
	if err != nil {
		log.Println("Error:", err)
	}
	fmt.Println("Login:", host, user, pw)
	Login(host, user, pw)
}

func register(host string, user string, pw string) {
	host, user, pw, err := areCredsComplete(host, user, pw)
	if err != nil {

	} else {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	fmt.Println("Register:", host, user, pw)
}

func areCredsComplete(host string, user string, pw string) (string, string, string, error) {
	// Check if entered credentials are complete
	var err error = nil
	host = strings.TrimSpace(host)
	user = strings.TrimSpace(user)
	pw = strings.TrimSpace(pw)

	if host == "" {
		host = HOST_DEFAULT
	}
	if !isHostnameValid(host) {
		mainWindow.Message("Please enter a valid hostname").WithError().Show()
	} else if user == "" {
		mainWindow.Message("Please enter an username").WithError().Show()
		err = errors.New("missing username")
	} else if pw == "" {
		mainWindow.Message("Please enter a password").WithError().Show()
		err = errors.New("missing username")
	}
	return host, user, pw, err
}

func isHostnameValid(host string) bool {
	match, _ := regexp.Match(RE_HOSTNAME, []byte(host))
	return match
}

func renderWindow() base.Widget {
	var user, pw, host string
	widget :=
		&goey.HBox{
			Children: []base.Widget{
				&goey.VBox{
					Children: []base.Widget{
						&goey.Label{Text: "Server address:"},
						&goey.TextInput{Value: host, Placeholder: "https://clipster.cc",
							OnChange: func(v string) {
								host = v
								println("server input ", v)
							}},
						&goey.Checkbox{Text: "Disable SSL cert check",
							OnChange: func(v bool) { println("check box input: ", v) }},
						&goey.Label{Text: "Username:"},
						&goey.TextInput{Value: user, Placeholder: "Enter username",
							OnChange: func(v string) {
								user = v
								println("username input ", v)
							}},
						&goey.Label{Text: "Password:"},
						&goey.TextInput{Value: pw, Placeholder: "Enter password", Password: true,
							OnChange: func(v string) {
								pw = v
								println("password input ", v)
							}},
						&goey.HBox{
							Children: []base.Widget{
								&goey.Button{Text: "Login", OnClick: func() {
									login(host, user, pw)
								}},
								&goey.Button{Text: "Register", OnClick: func() {
									register(host, user, pw)
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
