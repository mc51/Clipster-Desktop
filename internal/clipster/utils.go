// Utility functions used throughout package
package clipster

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/biessek/golang-ico"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/draw"

	"fyne.io/fyne/v2"

	"clipster/assets"
)

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
		img_decoded, _, err = image.Decode(bytes.NewReader(assets.NotFoundPNG))
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
	host = normalizeHost(host)
	user = strings.TrimSpace(user)
	pw = strings.TrimSpace(pw) // maybe space should be valid? but not at beginning or end?

	if !isHostnameValid(host) {
		err = errors.New("Please enter a valid hostname")
	} else if user == "" {
		err = errors.New("Please enter an username")
	} else if pw == "" {
		err = errors.New("Please enter a password")
	}
	return host, user, pw, err
}

// normalizeHost trims the address, empty means the default server
func normalizeHost(host string) string {
	host = strings.TrimRight(strings.TrimSpace(host), "/")
	if host == "" {
		return HOST_DEFAULT
	}
	return host
}

// isHostnameValid checks that host is a usable server address. The credentials are sent
// to it, so it must be https. Plain http is only accepted for the local machine.
// Credentials in the address (https://user:pw@host) are refused, as they would end up in logs
func isHostnameValid(host string) bool {
	u, err := url.Parse(host)
	if err != nil || u.Hostname() == "" || u.User != nil || u.RawQuery != "" ||
		u.Fragment != "" || u.Opaque != "" || strings.Contains(host, "?") || strings.Contains(host, "#") {
		return false
	}
	switch u.Scheme {
	case "https":
		return true
	case "http":
		return u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1"
	}
	return false
}

// login_flow checks for completeness of creds, creates hash from them and
// uses hash to authenticate against API endpoint. On success saves credentials to config.
// Blocks on the network, so do not call it from the GUI goroutine.
// Returns a message describing the result to be displayed to the user
func login_flow(host string, user string, pw string, pin string) (string, error) {
	host, user, pw, err := AreCredsComplete(host, user, pw)
	if err != nil {
		log.Println("Error:", err)
		return "", err
	}
	log.Println("Login: certificate pinned:", pin != "")
	debugf("Login: %s %s", host, user)

	hash_login := GetLoginHashFromPw(user, pw)
	if err := APILogin(host, user, hash_login, pin); err != nil {
		log.Println("Error:", err)
		return "", err
	}
	return saveCredentials(Config{host, user, hash_login, GetMsgHashFromPw(user, pw), pin},
		"Login successful")
}

// register_flow checks for completeness of creds, creates hash from them and
// uses hash to register at API endpoint. On success saves credentials to config.
// Blocks on the network, so do not call it from the GUI goroutine.
// Returns a message describing the result to be displayed to the user
func register_flow(host string, user string, pw string, pin string) (string, error) {
	host, user, pw, err := AreCredsComplete(host, user, pw)
	if err != nil {
		log.Println("Error:", err)
		return "", err
	}
	log.Println("Registration: certificate pinned:", pin != "")
	debugf("Registration: %s %s", host, user)

	hash_login := GetLoginHashFromPw(user, pw)
	if err := APIRegister(host, user, hash_login, pin); err != nil {
		log.Println("Error:", err)
		return "", err
	}
	return saveCredentials(Config{host, user, hash_login, GetMsgHashFromPw(user, pw), pin},
		"Registration successful")
}

// saveCredentials makes c the active config and writes it to disk.
// Returns the message to display to the user
func saveCredentials(c Config, msg string) (string, error) {
	setConf(c)
	if err := WriteConfigFile(c); err != nil {
		return "", fmt.Errorf("%s\nBut credentials could not be saved to config:\n%s\n%w",
			msg, CONFIG_FILEPATH, err)
	}
	log.Println("Ok:", msg)
	return msg + "\nCredentials saved to config:\n" + CONFIG_FILEPATH, nil
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
		fyne.Do(func() { GUI_AllClips(clips) })
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
	log.Println("Shared clip: format", format, "with length", len(clip))
	if format == "txt" { // images are huge base64 strings
		debugf("Shared clip: %s", clip)
	}
	ShowNotification("Clipster – Shared clip", MSG_NOTIFY_SHARED)
}

// processClipTextToImages decodes the image of a clip and adds its PNG bytes and the PNG
// bytes of a thumbnail to the clip. It does not touch the GUI, so it is safe to call from any goroutine
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
