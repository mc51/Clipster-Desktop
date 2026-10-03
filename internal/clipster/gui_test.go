package clipster

import (
	"image"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// newTestApp makes Fyne run headless for the duration of a test
func newTestApp(t *testing.T) {
	t.Helper()
	fyneApp = test.NewApp()
	conf = Config{}
	configWin = nil
	t.Cleanup(func() {
		fyneApp = nil
		configWin = nil
		test.NewApp() // resets the global state of Fyne
	})
}

// windowsWithTitle returns all open windows with the given title. The headless driver
// always owns one additional untitled window, which is why windows are not just counted
func windowsWithTitle(title string) []fyne.Window {
	var found []fyne.Window
	for _, w := range fyneApp.Driver().AllWindows() {
		if w.Title() == title {
			found = append(found, w)
		}
	}
	return found
}

// findButton returns the button with the given label somewhere below obj
func findButton(obj fyne.CanvasObject, label string) *widget.Button {
	if b, ok := obj.(*widget.Button); ok && b.Text == label {
		return b
	}
	var children []fyne.CanvasObject
	if c, ok := obj.(*fyne.Container); ok {
		children = c.Objects
	} else if w, ok := obj.(fyne.Widget); ok {
		children = test.WidgetRenderer(w).Objects()
	}
	for _, child := range children {
		if b := findButton(child, label); b != nil {
			return b
		}
	}
	return nil
}

func TestSanitizeNotification(t *testing.T) {
	tests := []struct{ in, want string }{
		{"plain text äöü", "plain text äöü"},
		{"$(calc.exe) $env:USERNAME", "＄(calc.exe) ＄env:USERNAME"},
		{`" ; calc ; "`, `' ; calc ; '`},
		{"“ ; calc ; ”", "' ; calc ; '"},
		{"`n back`tick", "'n back'tick"},
		{`C:\dir\file`, "C:/dir/file"},
		{"line1\nline2\r\n\ttab", "line1 line2   tab"},
	}
	for _, tt := range tests {
		got := sanitizeNotification(tt.in)
		if got != tt.want {
			t.Errorf("sanitizeNotification(%q) = %q, want %q", tt.in, got, tt.want)
		}
		if strings.ContainsAny(got, "\"$`\\\n") {
			t.Errorf("sanitizeNotification(%q) = %q still contains special characters", tt.in, got)
		}
	}
}

func TestTextRowHeight(t *testing.T) {
	if got := textRowHeight("short"); got != textRowMinHeight {
		t.Errorf("short text height = %v, want %v", got, textRowMinHeight)
	}
	if got := textRowHeight(strings.Repeat("x\n", 1000)); got != textRowMaxHeight {
		t.Errorf("long text height = %v, want %v", got, textRowMaxHeight)
	}
}

func TestConfigWindow(t *testing.T) {
	newTestApp(t)
	conf = Config{Server: "https://example.com", Username: "bob"}

	GUI_ConfigWindow()
	GUI_ConfigWindow() // must reuse the open window
	windows := windowsWithTitle("Clipster - Config")
	if len(windows) != 1 {
		t.Fatalf("got %d config windows, want 1", len(windows))
	}
	if findButton(windows[0].Content(), "Login") == nil ||
		findButton(windows[0].Content(), "Register") == nil ||
		findButton(windows[0].Content(), "Cancel") == nil {
		t.Error("config window misses a button")
	}
	if windows[0].Canvas().Capture() == nil {
		t.Error("config window can not be rendered")
	}

	test.Tap(findButton(windows[0].Content(), "Cancel"))
	if configWin != nil {
		t.Error("cancel should close the config window")
	}
}

func TestAllClipsWindow(t *testing.T) {
	newTestApp(t)
	thumb, err := ImageToBytes(image.NewRGBA(image.Rect(0, 0, 40, 20)))
	if err != nil {
		t.Fatal(err)
	}
	clips := []Clips{
		{Format: "txt", TextDecrypted: "hello\nworld"},
		{Format: "img", ThumbBytes: thumb, ImageBytes: thumb},
		{Format: "img"}, // no thumbnail: falls back to placeholder
	}

	GUI_AllClips(clips)
	windows := windowsWithTitle("Clipster - Your Clips")
	if len(windows) != 1 {
		t.Fatalf("got %d clips windows, want 1", len(windows))
	}
	w := windows[0]
	if w.Canvas().Capture() == nil {
		t.Error("clips window can not be rendered")
	}

	// nothing is selected, so copy and save must complain
	for _, label := range []string{"Copy", "Save as"} {
		btn := findButton(w.Content(), label)
		if btn == nil {
			t.Fatalf("clips window misses button %q", label)
		}
		test.Tap(btn)
		if w.Canvas().Overlays().Top() == nil {
			t.Errorf("%q without selection should show an error dialog", label)
		}
		w.Canvas().Overlays().Remove(w.Canvas().Overlays().Top())
	}
}
