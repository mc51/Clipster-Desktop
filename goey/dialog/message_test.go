package dialog

import (
	"fmt"
	"testing"

	"guitest/goey/loop"
)

func ExampleNewMessage() {
	// The following creates a modal dialog with a message.
	err := NewMessage("Some text for the body of the dialog box.").WithTitle("Example").WithInfo().Show()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func TestNewMessage(t *testing.T) {
	cases := []struct {
		build      func() error
		asyncEnter bool
		ok         bool
	}{
		{func() error {
			return NewMessage("Some text for the body of the dialog box.").WithTitle(t.Name()).WithInfo().Show()
		}, true, true},
		{func() error { return NewMessage("").Err() }, false, false},
		{func() error { return NewMessage("").Show() }, false, false},
		{func() error { return NewMessage("Some text...").WithTitle("").Err() }, false, false},
		{func() error { return NewMessage("Some text...").WithTitle("").Show() }, false, false},
	}

	init := func() error {
		for i, v := range cases {
			if v.asyncEnter {
				asyncKeyEnter()
			}

			err := v.build()
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
