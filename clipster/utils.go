package clipster

import (
	"encoding/base64"
	"errors"
	"image"
	"log"
	"os"
	"regexp"
	"strings"
)

func ReadIconAsB64FromFile(filename string) string {
	// ReadIconAsB64FromFile reads file content and returns it
	// as base64 encoded string
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	icon_b64 := base64.StdEncoding.EncodeToString(bytes)
	return icon_b64
}

func ReadIconAsBytesFromFile(filename string) []byte {
	// ReadIconAsBytesFromFile reads file content and returns it as bytes
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func ReadIconAsImageFromFile(filename string) (image.Image, error) {
	// ReadIconAsImageFromFile read file content and returns image type
	icon_b64 := ReadIconAsB64FromFile(filename)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(icon_b64))
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println("Error:", err)
	}
	return m, err
}

func AreCredsComplete(host string, user string, pw string) (string, string, string, error) {
	// AreCredsComplete check if entered credentials are complete
	var err error = nil
	host = strings.TrimSpace(host)
	user = strings.TrimSpace(user)
	pw = strings.TrimSpace(pw)

	if host == "" {
		host = HOST_DEFAULT
	}
	if !isHostnameValid(host) {
		err = errors.New("Please enter a valid hostname")
	} else if user == "" {
		err = errors.New("Please enter an username")
	} else if pw == "" {
		err = errors.New("Please enter a password")
	}
	return host, user, pw, err
}

func isHostnameValid(host string) bool {
	match, _ := regexp.Match(RE_HOSTNAME, []byte(host))
	return match
}
