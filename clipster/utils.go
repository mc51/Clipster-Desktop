// Utility functions used throughout package
package clipster

import (
	"bytes"
	"errors"
	"image"
	"log"
	"regexp"
	"strings"
)

func IconAsImageFromBytes(img []byte) (image.Image, error) {
	// IconAsImageFromBytes read image from bytes and returns image
	m, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		log.Panicln("Error:", err)
	}
	return m, err
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
