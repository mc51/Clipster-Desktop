// Deal with copy and paste to/from clipboard
package clipster

import (
	"log"

	"golang.design/x/clipboard"
)

var clip []byte

func GetClipboard() string {
	clip = clipboard.Read(clipboard.FmtText)
	log.Println("Get Clipboard:", string(clip))
	return string(clip)
}

func SetClipboard(clip string) {
	log.Println("Set Clipboard:", clip)
	clipboard.Write(clipboard.FmtText, []byte(clip))
}
