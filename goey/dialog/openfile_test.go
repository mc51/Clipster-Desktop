package dialog

import (
	"os"
	"path/filepath"
	"testing"

	"clipster/goey/loop"
)

func getwd(t *testing.T) string {
	path, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not determine the working directory, %s", err)
	}
	return path
}

func TestNewOpenFile(t *testing.T) {
	wd := getwd(t)

	cases := []struct {
		build    func() (string, error)
		asyncKey rune
		filename string
		ok       bool
	}{
		{func() (string, error) { return NewOpenFile().WithTitle(t.Name()).Show() }, '\x1b', "", true},
		{func() (string, error) { return "", NewOpenFile().WithTitle("").Err() }, 0, "", false},
		{func() (string, error) { return NewOpenFile().WithTitle("").Show() }, 0, "", false},
		{func() (string, error) { return NewOpenFile().WithFilename("./openfile_test.go").Show() }, '\n', filepath.Join(wd, "openfile_test.go"), true},
	}
	init := func() error {
		for i, v := range cases {
			if v.asyncKey == '\n' {
				asyncKeyEnter()
			} else if v.asyncKey == '\x1b' {
				asyncKeyEscape()
			}

			filename, err := v.build()
			if filename != v.filename {
				t.Errorf("Case %d, want %s, got %s", i, v.filename, filename)
			}
			if got := err == nil; got != v.ok {
				t.Errorf("Case %d,  want %v, got %v", i, v.ok, got)
				if err != nil {
					t.Logf("Error: %s", err)
				}
			}
		}

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
}
