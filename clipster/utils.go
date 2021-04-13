package clipster

import (
	"encoding/base64"
	"log"
	"os"
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
	// ReadIconAsBytesFromFile reads file content and returns it
	// as bytes
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
