// Deal with copy and paste to/from clipboard
package clipster

import (
	"encoding/base64"
	"log"

	"golang.design/x/clipboard"
)

// GetClipboard returns the current local clipboard and its content type.
// text is raw string, images are png as b64 standard encoded strings
func GetClipboard() (string, string) {
	clipBytes := clipboard.Read(clipboard.FmtText)
	if clipBytes == nil {
		clipBytes := clipboard.Read(clipboard.FmtImage)
		clip := base64.StdEncoding.EncodeToString(clipBytes)
		log.Println("Get Clipboard:", clip)
		return clip, "img"
	}
	clip := string(clipBytes)
	log.Println("Get Clipboard:", clip)
	return clip, "txt"
}

// SetClipboard moves clip content to local clipboard and shows notification.
// Deals with txt and img format
func SetClipboard(clip Clips) {
	log.Println(clip)
	if clip.Format == "img" {
		imgBytes, err := base64.StdEncoding.DecodeString(clip.TextDecrypted)
		if err != nil {
			log.Panicln("Error:", err)
		}
		clipboard.Write(clipboard.FmtImage, imgBytes)
		log.Println("Set Clipboard:", MSG_NOTIFY_GOT_IMAGE)
		ShowNotification("Clipster – Got new clip", MSG_NOTIFY_GOT_IMAGE)
	} else {
		clipboard.Write(clipboard.FmtText, []byte(clip.TextDecrypted))
		log.Println("Set Clipboard:", clip.TextDecrypted)
		ShowNotification("Clipster – Got new clip", clip.TextDecrypted)
	}
}
