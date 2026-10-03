package clipster

import (
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
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
	setConf(Config{Hash_msg: GetMsgHashFromPw("alice", "secret")})
	t.Cleanup(func() { setConf(Config{}) })

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

	setConf(Config{Hash_msg: GetMsgHashFromPw("alice", "other")})
	if _, err := Decrypt(tok); err == nil {
		t.Error("Decrypt with wrong key should fail")
	}
	setConf(Config{})
	if _, err := Encrypt("x"); err == nil {
		t.Error("Encrypt without key should fail")
	}
}

// Credentials are saved by the config window in the background, while flows started from
// the tray read them. Run with -race
func TestSaveCredentials(t *testing.T) {
	oldPath := CONFIG_FILEPATH
	CONFIG_FILEPATH = filepath.Join(t.TempDir(), CONFIG_FILENAME)
	t.Cleanup(func() {
		CONFIG_FILEPATH = oldPath
		setConf(Config{})
	})

	c := Config{"https://example.com", "alice", GetLoginHashFromPw("alice", "secret"),
		GetMsgHashFromPw("alice", "secret"), false}
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			if _, err := saveCredentials(c, "Login successful"); err != nil {
				t.Error(err)
			}
		})
		wg.Go(func() { _, _ = Encrypt("hello") })
	}
	wg.Wait()

	if getConf() != c {
		t.Errorf("active config = %+v, want %+v", getConf(), c)
	}
	info, err := os.Stat(CONFIG_FILEPATH)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Errorf("config file permissions = %v, want 0600", info.Mode().Perm())
	}
	if loaded, err := LoadConfigFromFile(); err != nil || loaded != c {
		t.Errorf("LoadConfigFromFile = %+v, %v, want %+v", loaded, err, c)
	}
}

func TestWriteConfigFileOverwrites(t *testing.T) {
	oldPath := CONFIG_FILEPATH
	CONFIG_FILEPATH = filepath.Join(t.TempDir(), "new folder", CONFIG_FILENAME)
	t.Cleanup(func() {
		CONFIG_FILEPATH = oldPath
		setConf(Config{})
	})

	if err := WriteConfigFile(Config{Server: "https://a.example", Username: "a", Hash_login: "l", Hash_msg: "m"}); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(CONFIG_FILEPATH, 0644); err != nil {
			t.Fatal(err)
		}
	}
	want := Config{"https://b.example", "b", "l2", "m2", true}
	if err := WriteConfigFile(want); err != nil {
		t.Fatal(err)
	}
	if got, err := LoadConfigFromFile(); err != nil || got != want {
		t.Errorf("LoadConfigFromFile = %+v, %v, want %+v", got, err, want)
	}
	if info, err := os.Stat(CONFIG_FILEPATH); err != nil {
		t.Fatal(err)
	} else if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Errorf("rewritten config file permissions = %v, want 0600", info.Mode().Perm())
	}
	if files, _ := os.ReadDir(filepath.Dir(CONFIG_FILEPATH)); len(files) != 1 {
		t.Errorf("config folder should only contain the config file, got %d files", len(files))
	}
}

func TestLoadConfigFromFileInvalid(t *testing.T) {
	dir := t.TempDir()
	oldPath := CONFIG_FILEPATH
	t.Cleanup(func() {
		CONFIG_FILEPATH = oldPath
		setConf(Config{})
	})
	files := map[string]string{
		"missing":    "",
		"not toml":   "server: https://clipster.cc\nusername: bob\n",
		"incomplete": "server = \"https://clipster.cc\"\nusername = \"bob\"\n",
	}
	for name, content := range files {
		CONFIG_FILEPATH = filepath.Join(dir, name+".toml")
		if name != "missing" {
			if err := os.WriteFile(CONFIG_FILEPATH, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := LoadConfigFromFile(); err == nil {
			t.Errorf("%s config should not load", name)
		}
		if getConf() != (Config{}) {
			t.Errorf("%s config must not become the active config", name)
		}
	}
}

func TestFindConfigFile(t *testing.T) {
	first, second, third := t.TempDir(), t.TempDir(), t.TempDir()
	dirs := []string{first, second, third}
	if got, want := findConfigFile(dirs), filepath.Join(first, CONFIG_FILENAME); got != want {
		t.Errorf("without any config file got %q, want %q", got, want)
	}
	// an old config.yaml does not count
	if err := os.WriteFile(filepath.Join(second, "config.yaml"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{third, second} {
		if err := os.WriteFile(filepath.Join(dir, CONFIG_FILENAME), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if got, want := findConfigFile(dirs), filepath.Join(second, CONFIG_FILENAME); got != want {
		t.Errorf("got %q, want first folder containing a config file %q", got, want)
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

func TestIsHostnameValid(t *testing.T) {
	tests := []struct {
		host string
		want bool
	}{
		{"https://example.com", true},
		{"https://example.com:8443/clipster", true},
		{"http://localhost:8000", true},
		{"http://127.0.0.1:8000", true},
		{"http://localhost", true},
		{"http://example.com", false},
		{"http://192.168.1.5:8000", false},
		{"http://evil.com/?a=://localhost:", false},
		{"http://evil.com/x://127.0.0.1:", false},
		{"http://localhost.evil.com", false},
		{"https://user:pw@example.com", false},
		{"https://user@example.com", false},
		{"https://example.com?x=1", false},
		{"https://example.com#x", false},
		{"https://", false},
		{"192.168.1.5", false},
		{"example.com", false},
		{"ftp://example.com", false},
		{"https://exa mple.com", false},
	}
	for _, tt := range tests {
		if got := isHostnameValid(tt.host); got != tt.want {
			t.Errorf("isHostnameValid(%q) = %v, want %v", tt.host, got, tt.want)
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
