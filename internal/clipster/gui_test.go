package clipster

import (
	"errors"
	"image"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// newTestApp makes Fyne run headless for the duration of a test
func newTestApp(t *testing.T) {
	t.Helper()
	fyneApp = test.NewApp()
	setConf(Config{})
	configWin = nil
	t.Cleanup(func() {
		fyneApp = nil
		configWin = nil
		setConf(Config{})
		loginFlow, registerFlow = login_flow, register_flow
		runInBackground = func(f func()) { go f() }
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

// walk calls f for obj and everything below it, until f returns false
func walk(obj fyne.CanvasObject, f func(fyne.CanvasObject) bool) bool {
	if !f(obj) {
		return false
	}
	var children []fyne.CanvasObject
	if c, ok := obj.(*fyne.Container); ok {
		children = c.Objects
	} else if w, ok := obj.(fyne.Widget); ok {
		children = test.WidgetRenderer(w).Objects()
	}
	for _, child := range children {
		if !walk(child, f) {
			return false
		}
	}
	return true
}

// findButton returns the button with the given label somewhere below obj
func findButton(obj fyne.CanvasObject, label string) *widget.Button {
	var found *widget.Button
	walk(obj, func(o fyne.CanvasObject) bool {
		if b, ok := o.(*widget.Button); ok && b.Text == label {
			found = b
		}
		return found == nil
	})
	return found
}

// findPasswordEntry returns the first password entry somewhere below obj
func findPasswordEntry(obj fyne.CanvasObject) *widget.Entry {
	var found *widget.Entry
	walk(obj, func(o fyne.CanvasObject) bool {
		if e, ok := o.(*widget.Entry); ok && e.Password {
			found = e
		}
		return found == nil
	})
	return found
}

// openConfigWindow opens the config window with flows that run synchronously
func openConfigWindow(t *testing.T) fyne.Window {
	t.Helper()
	runInBackground = func(f func()) { f() }
	setConf(Config{Server: "https://example.com", Username: "bob"})
	GUI_ConfigWindow()
	windows := windowsWithTitle("Clipster - Config")
	if len(windows) != 1 {
		t.Fatalf("got %d config windows, want 1", len(windows))
	}
	return windows[0]
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

func TestTextPreview(t *testing.T) {
	tests := []struct{ in, want string }{
		{"short", "short"},
		{"a\n\nb\n", "a\n \nb\n "},
		{"a\r\n\r\nb", "a\n \nb"},
	}
	for _, tt := range tests {
		if got := textPreview(tt.in); got != tt.want {
			t.Errorf("textPreview(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
	long := textPreview(strings.Repeat("ä", textPreviewMaxRunes+1))
	if !strings.HasSuffix(long, " [...]") || len([]rune(long)) > textPreviewMaxRunes+len(" [...]") {
		t.Errorf("long text is not truncated, got %d runes", len([]rune(long)))
	}
}

func TestTextRowHeight(t *testing.T) {
	newTestApp(t)
	if got := textRowHeight("short"); got != textRowMinHeight {
		t.Errorf("short text height = %v, want %v", got, textRowMinHeight)
	}
	if got := textRowHeight(strings.Repeat("x\n", 1000)); got != textRowMaxHeight {
		t.Errorf("long text height = %v, want %v", got, textRowMaxHeight)
	}
	medium := strings.Repeat("x\n", 7) + "x"
	if got, want := textRowHeight(medium), widget.NewLabel(medium).MinSize().Height; got != want {
		t.Errorf("8 lines height = %v, want height of label %v", got, want)
	}
}

func TestConfigWindow(t *testing.T) {
	newTestApp(t)
	w := openConfigWindow(t)
	GUI_ConfigWindow() // must reuse the open window
	if n := len(windowsWithTitle("Clipster - Config")); n != 1 {
		t.Fatalf("got %d config windows, want 1", n)
	}
	if findButton(w.Content(), "Login") == nil ||
		findButton(w.Content(), "Register") == nil ||
		findButton(w.Content(), "Cancel") == nil {
		t.Error("config window misses a button")
	}
	if w.Canvas().Capture() == nil {
		t.Error("config window can not be rendered")
	}

	test.Tap(findButton(w.Content(), "Cancel"))
	if configWin != nil {
		t.Error("cancel should close the config window")
	}
}

func TestConfigWindowLogin(t *testing.T) {
	newTestApp(t)
	w := openConfigWindow(t)
	loginBtn, registerBtn := findButton(w.Content(), "Login"), findButton(w.Content(), "Register")
	test.Type(findPasswordEntry(w.Content()), "secret")

	loginFlow = func(host, user, pw string, noSSL bool) (string, error) {
		if host != "https://example.com" || user != "bob" || pw != "secret" || noSSL {
			t.Errorf("login with %q %q %q %v", host, user, pw, noSSL)
		}
		if !loginBtn.Disabled() || !registerBtn.Disabled() {
			t.Error("buttons must be disabled while logging in")
		}
		return "Login successful", nil
	}
	test.Tap(loginBtn)

	info := w.Canvas().Overlays().Top()
	if info == nil {
		t.Fatal("login should show an info dialog")
	}
	if configWin != w {
		t.Fatal("config window should stay open until the info dialog is closed")
	}
	test.Tap(findButton(info, "OK"))
	if configWin != nil {
		t.Error("closing the info dialog should close the config window")
	}
}

func TestConfigWindowRegisterError(t *testing.T) {
	newTestApp(t)
	w := openConfigWindow(t)
	loginBtn, registerBtn := findButton(w.Content(), "Login"), findButton(w.Content(), "Register")

	registerFlow = func(string, string, string, bool) (string, error) {
		return "", errors.New("registration failed: user exists")
	}
	test.Tap(registerBtn)

	if w.Canvas().Overlays().Top() == nil {
		t.Error("failed registration should show an error dialog")
	}
	if configWin != w {
		t.Error("config window should stay open after an error")
	}
	if loginBtn.Disabled() || registerBtn.Disabled() {
		t.Error("buttons must be enabled again after an error")
	}
}

// When the user closes the config window while the request is running,
// the result must still reach them
func TestConfigWindowClosedDuringRequest(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		err  error
		want *fyne.Notification
	}{
		{"success", "Login successful", nil, fyne.NewNotification("Clipster", "Login successful")},
		{"error", "", errors.New("login failed"), fyne.NewNotification("Clipster - Error", "login failed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTestApp(t)
			w := openConfigWindow(t)
			loginFlow = func(string, string, string, bool) (string, error) {
				w.Close()
				return tt.msg, tt.err
			}
			test.AssertNotificationSent(t, tt.want, func() {
				test.Tap(findButton(w.Content(), "Login"))
			})
			if configWin != nil {
				t.Error("config window should be closed")
			}
		})
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

// Rows of the clips list do not clip their content, so long texts must be cut to the row
func TestAllClipsWindowLongTextStaysInRow(t *testing.T) {
	texts := map[string]string{
		"many lines":       strings.Repeat("line of text\n", 60),
		"many empty lines": "start" + strings.Repeat("\n", 60) + "end",
		"one long line":    strings.Repeat("word ", 2000),
		"mixed":            strings.Repeat(strings.Repeat("word ", 30)+"\n\n", 20),
	}
	for name, text := range texts {
		t.Run(name, func(t *testing.T) {
			newTestApp(t)
			GUI_AllClips([]Clips{{Format: "txt", TextDecrypted: text}, {Format: "txt", TextDecrypted: "next"}})
			w := windowsWithTitle("Clipster - Your Clips")[0]
			w.Resize(fyne.NewSize(700, 900))
			w.Canvas().Capture() // lays out the list

			checked := 0
			walk(w.Content(), func(o fyne.CanvasObject) bool {
				rt, ok := o.(*widget.RichText)
				if !ok || !rt.Visible() {
					return true
				}
				for _, obj := range test.WidgetRenderer(rt).Objects() {
					if txt, ok := obj.(*canvas.Text); ok && txt.Visible() &&
						txt.Position().Y+txt.MinSize().Height > rt.Size().Height+0.5 {
						t.Errorf("text %q ends at %v, below its row of height %v",
							txt.Text, txt.Position().Y+txt.MinSize().Height, rt.Size().Height)
						return false
					}
				}
				checked++
				return true
			})
			if checked == 0 {
				t.Fatal("found no text in the clips list")
			}
		})
	}
}
