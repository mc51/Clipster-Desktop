// Utility functions used throughout package
package clipster

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"log"
	"regexp"
	"strings"
)

// BytesToImage reads image from bytes and returns image.Image
func BytesToImage(img []byte) (image.Image, error) {
	m, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		log.Panicln("Error:", err)
	}
	return m, err
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

func AreCredsComplete(host string, user string, pw string) (string, string, string, error) {
	// AreCredsComplete check if entered credentials are complete and hostname is valid
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

func isHostnameValid(host string) bool {
	// isHostnameValid checks hostname against some regex for basic validity
	match, _ := regexp.Match(RE_HOSTNAME, []byte(host))
	return match
}
