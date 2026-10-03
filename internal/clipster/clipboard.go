// Deal with copy and paste to/from clipboard
package clipster

import (
	"context"
	"encoding/base64"
	"errors"
	"log"

	"golang.design/x/clipboard"
)

// GetClipboard returns the current local clipboard and its content type.
// text is raw string, images are png as b64 standard encoded strings
func GetClipboard() (string, string, error) {
	ctx := context.Background()
	clipBytes, err := clipboard.Read(ctx, clipboard.FmtText)
	if err == nil && len(clipBytes) > 0 {
		log.Println("Get Clipboard: text with length", len(clipBytes))
		return string(clipBytes), "txt", nil
	}
	if err != nil && !errors.Is(err, clipboard.ErrNoData) {
		return "", "", err
	}

	clipBytes, err = clipboard.Read(ctx, clipboard.FmtImage)
	if err == nil && len(clipBytes) > 0 {
		log.Println("Get Clipboard: image with size", len(clipBytes))
		return base64.StdEncoding.EncodeToString(clipBytes), "img", nil
	}
	if err != nil && !errors.Is(err, clipboard.ErrNoData) {
		return "", "", err
	}
	return "", "", errors.New("clipboard is empty")
}

// SetClipboard moves clip content to local clipboard and shows notification.
// Deals with txt and img format
func SetClipboard(clip Clips) {
	ctx := context.Background()
	if clip.Format == "img" {
		if _, err := clipboard.Write(ctx, clipboard.FmtImage, clip.ImageBytes); err != nil {
			log.Println("Error: set clipboard", err)
			ShowNotification("Clipster - Error", err.Error())
			return
		}
		log.Println("Set Clipboard: image with size", len(clip.ImageBytes))
	} else {
		if _, err := clipboard.Write(ctx, clipboard.FmtText, []byte(clip.TextDecrypted)); err != nil {
			log.Println("Error: set clipboard", err)
			ShowNotification("Clipster - Error", err.Error())
			return
		}
		log.Println("Set Clipboard: text with length", len(clip.TextDecrypted))
		debugf("Set Clipboard: %s", clip.TextDecrypted)
	}
	ShowNotification("Clipster – Got new clip", MSG_NOTIFY_GOT_CLIP)
}
