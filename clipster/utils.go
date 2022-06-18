// Utility functions used throughout package
package clipster

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	_ "github.com/biessek/golang-ico"
	_ "golang.org/x/image/bmp"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/nfnt/resize"
)

// BytesToPixbuf takes Image in bytes and returns gdk.Pixbuf representation
func BytesToPixbuf(img []byte) *gdk.Pixbuf {
	i, err := gdk.PixbufNewFromBytesOnly(img)
	if err != nil {
		log.Println("Could not create icon", err)
	}
	return i
}

// BytesToImage reads bytes and returns image.Image. If bytes are not a valid Image
// return a default "file not found" Image
func BytesToImage(img []byte) (image.Image, error) {

	mimeType := http.DetectContentType(img)
	log.Printf("BytesToImage mimeType: %s", mimeType)

	img_decoded, format, err := image.Decode(bytes.NewReader(img))
	log.Printf("BytesToImage Decode Format: %s", format)
	if err != nil {
		log.Println("Error BytesToImage:", err)
		log.Println("Returning 'missing file' image instead", format)
		img_decoded, _, err = image.Decode(bytes.NewReader(PNG_BYTES_IMAGE_NOTFOUND))
	}
	return img_decoded, err
}

// ImageToBytes reads image and returns bytes
func ImageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		log.Println("Error Encode:", err)
	}
	img_bytes := buf.Bytes()
	return img_bytes, err
}

// B64ToImage converts b64 encoded string of an image to image.Image if it fails
// nil is returned
func B64ToImage(img string) (image.Image, error) {
	img_bytes, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		log.Println("Error DecodeString:", err)
	}
	image, err := BytesToImage(img_bytes)
	if err != nil {
		log.Println("Error BytesToImage:", err)
	}
	return image, err
}

// AreCredsComplete checks if entered credentials are complete and hostname is valid
func AreCredsComplete(host string, user string, pw string) (string, string, string, error) {
	var err error = nil
	host = strings.TrimSpace(host)
	user = strings.TrimSpace(user)
	pw = strings.TrimSpace(pw) // maybe space should be valid? but not at beginning or end?

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

// DownloadClipsFlow downloads all clips from API, unencrypts text and
// displays result. If last_only == True, instead last clip is moved to clipboard
func DownloadClipsFlow(last_only bool) {
	clips, err := APIDownloadAllClips()
	if err != nil {
		ShowNotification("Clipster - Error", err.Error())
		log.Println("Error:", err)
		return
	}
	log.Printf("Clips: %+v", clips)

	for i := range clips {
		clips[i].TextDecrypted = Decrypt(clips[i].Text)
		if clips[i].Format == "img" {
			clips[i] = processClipTextToImages(clips[i])
		}
	}

	if last_only {
		SetClipboard(clips[len(clips)-1])
	} else {
		DoGUI(func() { GUI_AllClips(clips) })
	}
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

// processClipTextToImages creates a gtk.Image Thumbnail from the original clip Image.
// Also creates a bytes representation of the Image. Adds all that to the Clip
func processClipTextToImages(clip Clips) Clips {

	img, err := B64ToImage(clip.TextDecrypted)
	img_thumb := resize.Thumbnail(THUMBNAIL_WIDTH, THUMBNAIL_HEIGHT, img, resize.NearestNeighbor)
	img_thumb_bytes, err := ImageToBytes(img_thumb)
	img_thumb_pixbuf, err := gdk.PixbufNewFromBytesOnly(img_thumb_bytes)

	clip.ImageBytes, err = ImageToBytes(img)
	clip.GtkThumb, err = gtk.ImageNewFromPixbuf(img_thumb_pixbuf)

	if err != nil {
		log.Println("Error processClipTextToImages:", err)
	}

	return clip
}

// ImageToDisk shows file saving dialog and saves image file at chosen path
func ImageToDisk(clip Clips) {

	path := GUI_FileChooserDialog()
	if path != "" {
		if clip.Format != "img" {
			GUI_DialogError("You can only save images to file!")
			return
		}
		err := ioutil.WriteFile(path, clip.ImageBytes, 0644)
		if err != nil {
			log.Println("Error: saving file", err)
			GUI_DialogError("Error saving file: " + path + "\n" + err.Error())
		}
		log.Println("Saved file: " + path)
		GUI_DialogInfo("File saved: " + path)
	}
}
