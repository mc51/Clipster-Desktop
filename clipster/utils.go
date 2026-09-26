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
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	_ "github.com/biessek/golang-ico"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/draw"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
)

var reHostname = regexp.MustCompile(RE_HOSTNAME)

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
		log.Println("Returning 'missing file' image instead")
		img_decoded, _, err = image.Decode(bytes.NewReader(PNG_BYTES_IMAGE_NOTFOUND))
	}
	return img_decoded, err
}

// ImageToBytes reads image and returns bytes
func ImageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		log.Println("Error Encode:", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// B64ToImage converts b64 encoded string of an image to image.Image. If the
// string is no valid image, a "file not found" image is returned instead
func B64ToImage(img string) (image.Image, error) {
	img_bytes, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		log.Println("Error DecodeString:", err)
	}
	return BytesToImage(img_bytes)
}

// Thumbnail scales img down to fit into maxWidth x maxHeight keeping the aspect ratio.
// Images that already fit are returned unchanged
func Thumbnail(maxWidth int, maxHeight int, img image.Image) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= maxWidth && h <= maxHeight {
		return img
	}
	if w*maxHeight > h*maxWidth {
		h = max(1, h*maxWidth/w)
		w = maxWidth
	} else {
		w = max(1, w*maxHeight/h)
		h = maxHeight
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, b, draw.Src, nil)
	return dst
}

// truncate shortens s to at most n runes and marks it as shortened
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + " [...]"
}

// AreCredsComplete checks if entered credentials are complete and hostname is valid
func AreCredsComplete(host string, user string, pw string) (string, string, string, error) {
	var err error = nil
	host = strings.TrimRight(strings.TrimSpace(host), "/")
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
	return reHostname.MatchString(host)
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
	log.Println("Login:", host, user, ssl_disable)

	hash_login := GetLoginHashFromPw(user, pw)
	// TODO: This is blocking. Goroutine?
	if err := APILogin(host, user, hash_login, ssl_disable); err != nil {
		log.Println("Error:", err)
		GUI_DialogError("Error: " + err.Error())
		return err
	}
	return saveCredentials(Config{host, user, hash_login, GetMsgHashFromPw(user, pw), ssl_disable},
		"Login successful")
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
	log.Println("Registration:", host, user, ssl_disable)

	hash_login := GetLoginHashFromPw(user, pw)
	// TODO: This is blocking. Goroutine?
	if err := APIRegister(host, user, hash_login, ssl_disable); err != nil {
		log.Println("Error:", err)
		GUI_DialogError("Error: " + err.Error())
		return err
	}
	return saveCredentials(Config{host, user, hash_login, GetMsgHashFromPw(user, pw), ssl_disable},
		"Registration successful")
}

// saveCredentials makes c the active config, writes it to disk and displays the result
func saveCredentials(c Config, msg string) error {
	conf = c
	if err := WriteConfigFile(conf); err != nil {
		GUI_DialogError(msg + "\nBut credentials could not be saved to config:\n" +
			CONFIG_FILEPATH + "\n" + err.Error())
		return err
	}
	log.Println("Ok:", msg)
	GUI_DialogInfo(msg + "\nCredentials saved to config:\n" + CONFIG_FILEPATH)
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
	if len(clips) == 0 {
		ShowNotification("Clipster", "There are no shared clips yet")
		return
	}
	if last_only {
		clips = clips[len(clips)-1:]
	}

	for i := range clips {
		clips[i].TextDecrypted, err = Decrypt(clips[i].Text)
		if err != nil {
			log.Println("Error:", err)
			if last_only {
				ShowNotification("Clipster - Error", err.Error())
				return
			}
			clips[i].Format = "txt"
			clips[i].TextDecrypted = "[Error: " + err.Error() + "]"
			continue
		}
		if clips[i].Format == "img" {
			clips[i] = processClipTextToImages(clips[i])
		}
	}

	if last_only {
		SetClipboard(clips[0])
	} else {
		DoGUI(func() { GUI_AllClips(clips) })
	}
}

// ShareClipFlow gets current clipboard value, encrypts it
// uploads it to server and shows notification
func ShareClipFlow() {
	log.Println("ShareClipFlow")
	clip, format, err := GetClipboard()
	if err != nil {
		ShowNotification("Clipster - Error", "Could not share clip: "+err.Error())
		log.Println("Error:", err)
		return
	}
	clip_encrypted, err := Encrypt(clip)
	if err != nil {
		ShowNotification("Clipster - Error", err.Error())
		log.Println("Error:", err)
		return
	}
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

// processClipTextToImages decodes the image of a clip and adds its PNG bytes and the PNG
// bytes of a thumbnail to the clip. It does not use GTK, so it is safe to call from any goroutine
func processClipTextToImages(clip Clips) Clips {
	img, err := B64ToImage(clip.TextDecrypted)
	if err != nil {
		log.Println("Error processClipTextToImages:", err)
		return clip
	}
	if clip.ImageBytes, err = ImageToBytes(img); err != nil {
		log.Println("Error processClipTextToImages:", err)
	}
	thumb := Thumbnail(THUMBNAIL_WIDTH, THUMBNAIL_HEIGHT, img)
	if clip.ThumbBytes, err = ImageToBytes(thumb); err != nil {
		log.Println("Error processClipTextToImages:", err)
	}
	return clip
}

// ImageToDisk shows file saving dialog and saves image file at chosen path
func ImageToDisk(clip Clips) {
	if clip.Format != "img" {
		GUI_DialogError("You can only save images to file!")
		return
	}
	path := GUI_FileChooserDialog()
	if path == "" {
		return
	}
	if err := os.WriteFile(path, clip.ImageBytes, 0644); err != nil {
		log.Println("Error: saving file", err)
		GUI_DialogError("Error saving file: " + path + "\n" + err.Error())
		return
	}
	log.Println("Saved file: " + path)
	GUI_DialogInfo("File saved: " + path)
}
