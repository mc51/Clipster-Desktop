package clipster

import (
	"image"
	"strings"
	"testing"
)

// Keys must stay identical to the other Clipster clients
func TestDeriveKey(t *testing.T) {
	if got, want := GetLoginHashFromPw("alice", "secret"), "X6h-rcfgNSrD5ieTJzuuYWelPyhhOw5c2hYMQ4b4Ego="; got != want {
		t.Errorf("GetLoginHashFromPw = %q, want %q", got, want)
	}
	if got, want := GetMsgHashFromPw("alice", "secret"), "CBkssJrp5I6nCP6mt-YOIOU3sf4sOJlV7msooo_JJoQ="; got != want {
		t.Errorf("GetMsgHashFromPw = %q, want %q", got, want)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	conf = Config{Hash_msg: GetMsgHashFromPw("alice", "secret")}
	t.Cleanup(func() { conf = Config{} })

	tok, err := Encrypt("hello clipster")
	if err != nil {
		t.Fatal(err)
	}
	msg, err := Decrypt(tok)
	if err != nil {
		t.Fatal(err)
	}
	if msg != "hello clipster" {
		t.Errorf("Decrypt = %q", msg)
	}

	conf.Hash_msg = GetMsgHashFromPw("alice", "other")
	if _, err := Decrypt(tok); err == nil {
		t.Error("Decrypt with wrong key should fail")
	}
	conf.Hash_msg = ""
	if _, err := Encrypt("x"); err == nil {
		t.Error("Encrypt without key should fail")
	}
}

func TestAreCredsComplete(t *testing.T) {
	tests := []struct {
		host, user, pw string
		wantHost       string
		wantErr        bool
	}{
		{"", "u", "p", HOST_DEFAULT, false},
		{" https://example.com/ ", "u", "p", "https://example.com", false},
		{"http://localhost:8000", "u", "p", "http://localhost:8000", false},
		{"http://example.com", "u", "p", "http://example.com", true},
		{"", " ", "p", HOST_DEFAULT, true},
		{"", "u", "", HOST_DEFAULT, true},
	}
	for _, tt := range tests {
		host, _, _, err := AreCredsComplete(tt.host, tt.user, tt.pw)
		if host != tt.wantHost || (err != nil) != tt.wantErr {
			t.Errorf("AreCredsComplete(%q, %q, %q) = %q, %v", tt.host, tt.user, tt.pw, host, err)
		}
	}
}

func TestTruncate(t *testing.T) {
	if got := truncate("short", 10); got != "short" {
		t.Errorf("truncate = %q", got)
	}
	got := truncate(strings.Repeat("ä", 5), 3)
	if got != "äää [...]" {
		t.Errorf("truncate must not split runes, got %q", got)
	}
}

func TestThumbnail(t *testing.T) {
	tests := []struct{ w, h, wantW, wantH int }{
		{100, 50, 100, 50},
		{400, 200, 200, 100},
		{200, 800, 50, 200},
		{2000, 1, 200, 1},
	}
	for _, tt := range tests {
		img := image.NewRGBA(image.Rect(0, 0, tt.w, tt.h))
		b := Thumbnail(THUMBNAIL_WIDTH, THUMBNAIL_HEIGHT, img).Bounds()
		if b.Dx() != tt.wantW || b.Dy() != tt.wantH {
			t.Errorf("Thumbnail(%dx%d) = %dx%d, want %dx%d", tt.w, tt.h, b.Dx(), b.Dy(), tt.wantW, tt.wantH)
		}
	}
}

func TestProcessClipTextToImagesInvalid(t *testing.T) {
	clip := processClipTextToImages(Clips{Format: "img", TextDecrypted: "not an image"})
	if len(clip.ImageBytes) == 0 || len(clip.ThumbBytes) == 0 {
		t.Error("invalid image should fall back to placeholder image")
	}
}

func TestDesktopEntryExec(t *testing.T) {
	if got, want := desktopEntryExec(`/opt/my apps/clip$ter`), `"/opt/my apps/clip\\$ter"`; got != want {
		t.Errorf("desktopEntryExec = %q, want %q", got, want)
	}
}
