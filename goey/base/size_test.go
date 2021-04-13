package base

import (
	"fmt"
	"image"
	"testing"
)

func ExampleFromPixels() {
	// Most code should not need to worry about setting the DPI.  Windows will
	// ensure that the DPI is set.
	DPI = image.Point{96, 96}

	size := FromPixels(48, 96+96)
	fmt.Printf("The size is %s.\n", size.String())

	// Output:
	// The size is (48:00x192:00).
}

func TestSize(t *testing.T) {
	cases := []struct {
		in     Size
		isZero bool
		out    string
	}{
		{Size{}, true, "(0:00x0:00)"},
		{Size{1, 2}, false, "(0:01x0:02)"},
		{Size{1 * DIP, 2 * DIP}, false, "(1:00x2:00)"},
	}

	for i, v := range cases {
		if out := v.in.IsZero(); out != v.isZero {
			t.Errorf("Case %d:  Failed predicate IsZero, got %v, want %v", i, out, v.isZero)
		}
		if out := v.in.String(); out != v.out {
			t.Errorf("Case %d:  Failed method String, got %v, want %v", i, out, v.out)
		}
	}
}
