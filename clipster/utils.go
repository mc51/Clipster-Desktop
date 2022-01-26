// Utility functions used throughout package
package clipster

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/png"
	"log"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/glib"
)

// BytesToImage reads image from bytes and returns image.Image
func BytesToImage(img []byte) (image.Image, error) {
	m, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		log.Panicln("Error:", err)
	}
	return m, err
}

// ImageToBytes reads image and returns bytes
func ImageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	img_bytes := buf.Bytes()
	return img_bytes, err
}

// B64ToImage converts b64 encoded string of an image to image.Image
func B64ToImage(img string) (image.Image, error) {
	img_bytes, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		log.Panicln("Error:", err)
	}
	image, err := BytesToImage(img_bytes)
	if err != nil {
		log.Panicln("Error:", err)
	}
	return image, err
}

// AreCredsComplete checks if entered credentials are complete and hostname is valid
func AreCredsComplete(host string, user string, pw string) (string, string, string, error) {
	var err error = nil
	host = strings.TrimSpace(host)
	user = strings.TrimSpace(user)
	pw = strings.TrimSpace(pw)

	if host == "" {
		host = HOST_DEFAULT
	}
	if !isHostnameValid(host) {
		err = errors.New(" Please enter a valid hostname")
	} else if user == "" {
		err = errors.New(" Please enter an username")
	} else if pw == "" {
		err = errors.New(" Please enter a password")
	}
	return host, user, pw, err
}

// isHostnameValid checks hostname against some regex for basic validity
func isHostnameValid(host string) bool {
	match, _ := regexp.Match(RE_HOSTNAME, []byte(host))
	return match
}

// DoGUI adds function to be run on GTK Main loop / main thread
func DoGUI(action func()) {
	// Native GTK is not thread safe, and thus, gotk3's GTK bindings may not
	// be used from other goroutines.  Instead, glib.IdleAdd() must be used
	// to add a function to run in the GTK main loop when it is in an idle
	// state. See:
	// https://github.com/gotk3/gotk3-examples/blob/master/gtk-examples/goroutines/goroutines.go
	glib.IdleAdd(action)
}

// login_flow check for completeness of creds, creates hash from them,
// uses hash to authenticate against API endpoint, displays Message box with the result.
// On success saves credentials to config
func login_flow(host string, user string, pw string, ssl_disable bool) error {
	host, user, pw, err := AreCredsComplete(host, user, pw)
	if err != nil {
		GUI_DialogError("Error: " + err.Error())
		log.Println("Error:", err)
		return err
	}
	// TODO: Remove all cleartext pws from logs?
	log.Println("Login:", host, user, pw, ssl_disable)

	hash_login := GetLoginHashFromPw(user, pw)
	// TODO: This is blocking. Goroutine?
	if err := APILogin(host, user, hash_login, ssl_disable); err != nil {
		log.Println("Error:", err)
		GUI_DialogError("Error: " + err.Error())
		return err
	}
	hash_msg := GetMsgHashFromPw(user, pw)
	conf = Config{host, user, hash_login, hash_msg, ssl_disable}
	WriteConfigFile(conf)
	log.Println("Ok: login workflow completed")
	GUI_DialogInfo("Login successfull\nCredentials saved to config:\n" +
		CONFIG_FILEPATH)
	return nil
}

// register_flow check for completeness of creds, creates hash from them,
// uses hash to register at API endpoint, displays Message box with the result.
// On success saves credentials to config
func register_flow(host string, user string, pw string, ssl_disable bool) error {
	host, user, pw, err := AreCredsComplete(host, user, pw)
	if err != nil {
		GUI_DialogError("Error: " + err.Error())
		log.Println("Error:", err)
		return err
	}
	// TODO: Remove all cleartext pws from logs?
	log.Println("Registration:", host, user, pw, ssl_disable)

	hash_login := GetLoginHashFromPw(user, pw)
	// TODO: This is blocking. Goroutine?
	if err := APIRegister(host, user, hash_login, ssl_disable); err != nil {
		log.Println("Error:", err)
		GUI_DialogError("Error: " + err.Error())
		return err
	}
	hash_msg := GetMsgHashFromPw(user, pw)
	conf = Config{host, user, hash_login, hash_msg, ssl_disable}
	WriteConfigFile(conf)
	log.Println("Ok: Registration flow completed")

	GUI_DialogInfo("Registration successfull\nCredentials saved to config:\n" +
		CONFIG_FILEPATH)
	return nil
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
	DoGUI(func() { GUI_AllClips(clips) })
}

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
