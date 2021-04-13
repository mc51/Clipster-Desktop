package base

import (
	"fmt"
	"image"
	"testing"
)

func ExampleLength() {
	// Since there are 96 device-independent pixels per inch, and 6 picas
	// per inch, the following two lengths should be equal.
	length1 := 96 * DIP
	length2 := 6 * PC

	if length1 == length2 {
		fmt.Printf("All is OK with the world.")
	} else {
		fmt.Printf("This should not happen, unless there is a rounding error.")
	}

	// Output:
	// All is OK with the world.
}

func ExampleLength_Scale() {
	// There are 96 DIP in an inch, and 6 pica in a inch, so the following
	// should work.

	if length := (1 * DIP).Scale(96, 6); length == (1 * PC) {
		fmt.Printf("The ratio of pica to DIP is 96 to 6.")
	}

	// Output:
	// The ratio of pica to DIP is 96 to 6.
}

func ExampleLength_String() {
	fmt.Printf("Converting:  1pt is equal to %sdip\n", 1*PT)
	fmt.Printf("Converting:  1pt is equal to %1.2fdip\n", (1 * PT).DIP())
	fmt.Printf("Converting:  1pc is equal to %1.1fdip\n", (1 * PC).DIP())

	// Output:
	// Converting:  1pt is equal to 1:21dip
	// Converting:  1pt is equal to 1.33dip
	// Converting:  1pc is equal to 16.0dip
}

func ExampleRectangle() {
	r := Rectangle{Point{10 * DIP, 20 * DIP}, Point{90 * DIP, 80 * DIP}}

	fmt.Printf("Rectangle %s has dimensions %.0fdip by %.0fdip.",
		r, r.Dx().DIP(), r.Dy().DIP(),
	)

	// Output:
	// Rectangle (10:00,20:00)-(90:00,80:00) has dimensions 80dip by 60dip.
}

func ExampleRectangle_Pixels() {
	// The following line is for the example only, and should not appear in
	// user code, as the platform-specific code should update the DPI based
	// on the system.  However, for the purpose of this example, set a known
	// DPI.
	DPI = image.Point{2 * 96, 2 * 96}

	// Construct an example rectangle.
	r := Rectangle{Point{10 * DIP, 20 * DIP}, Point{90 * DIP, 80 * DIP}}
	rpx := r.Pixels()

	fmt.Printf("Rectangle %s when translated to pixels is %s.", r, rpx)

	// Output:
	// Rectangle (10:00,20:00)-(90:00,80:00) when translated to pixels is (20,40)-(180,160).
}

func TestFromPixels(t *testing.T) {
	cases := []struct {
		dpix, dpiy       int
		pixelsx, pixelsy int
		lengthx, lengthy Length
	}{
		{96, 96, 2, 3, 2 * DIP, 3 * DIP},
		{96, 96 * 3 / 2, 2, 3, 2 * DIP, 2 * DIP},
	}

	for i, v := range cases {
		DPI = image.Point{v.dpix, v.dpiy}
		if got := FromPixelsX(v.pixelsx); got != v.lengthx {
			t.Errorf("Unexpected conversion in FromPixelsX on case %d, got %v, want %v", i, got, v.lengthx)
		}
		if got := FromPixelsY(v.pixelsy); got != v.lengthy {
			t.Errorf("Unexpected conversion in FromPixelsY on case %d, got %v, want %v", i, got, v.lengthy)
		}
	}
}

func TestLength(t *testing.T) {
	if rt := (1 * DIP).DIP(); rt != 1 {
		t.Errorf("Unexpected round-trip for Length, %v =/= %v", rt, 1)
	}
	if rt := (1 * PT).PT(); rt != 1 {
		t.Errorf("Unexpected round-trip for PT,  %v =/= %v", rt, 1)
	}
	if rt := (1 * PC).PC(); rt != 1 {
		t.Errorf("Unexpected round-trip for PC,  %v =/= %v", rt, 1)
	}
	if rt := (1 * Inch).Inch(); rt != 1 {
		t.Errorf("Unexpected round-trip for inch,  %v =/= %v", rt, 1)
	}
	if rt := (1 * PT) * (1 << 6) / (1 * DIP); rt != 96*(1<<6)/72 {
		t.Errorf("Unexpected ratio between DIP and PT, %v =/= %v", rt, 96*(1<<6)/72)
	}
	if rt := (1 * PC) * (1 << 6) / (1 * DIP); rt != 96*(1<<6)/6 {
		t.Errorf("Unexpected ratio between DIP and PC, %v =/= %v", rt, 96*(1<<6)/72)
	}
	if rt := (1 * Inch) * (1 << 6) / (1 * DIP); rt != 96*(1<<6) {
		t.Errorf("Unexpected ratio between DIP and inch, %v =/= %v", rt, 96*(1<<6))
	}
}

func TestLength_Clamp(t *testing.T) {
	cases := []struct {
		in       Length
		min, max Length
		out      Length
	}{
		{10 * DIP, 0 * DIP, 20 * DIP, 10 * DIP},
		{30 * DIP, 0 * DIP, 20 * DIP, 20 * DIP},
		{-10 * DIP, 0 * DIP, 20 * DIP, 0 * DIP},
		{10 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{30 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{-10 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{10 * DIP, 20 * DIP, 0 * DIP, 20 * DIP},
		{30 * DIP, 20 * DIP, 0 * DIP, 20 * DIP},
		{-10 * DIP, 20 * DIP, 0 * DIP, 20 * DIP},
	}

	for i, v := range cases {
		if out := v.in.Clamp(v.min, v.max); out != v.out {
			t.Errorf("Error in case %d, want %s, got %s", i, v.out, out)
		}
	}
}

func TestPoint(t *testing.T) {
	cases := []struct {
		a, b Point
		add  Point
		sub  Point
	}{
		{Point{}, Point{1 * DIP, 2 * DIP}, Point{1 * DIP, 2 * DIP}, Point{-1 * DIP, -2 * DIP}},
		{Point{1 * DIP, 2 * DIP}, Point{}, Point{1 * DIP, 2 * DIP}, Point{1 * DIP, 2 * DIP}},
		{Point{3 * DIP, 5 * DIP}, Point{7 * DIP, 11 * DIP}, Point{10 * DIP, 16 * DIP}, Point{-4 * DIP, -6 * DIP}},
	}

	for i, v := range cases {
		if out := v.a.Add(v.b); out != v.add {
			t.Errorf("Error in case %d, want %s, got %s", i, v.add, out)
		}
		if out := v.a.Sub(v.b); out != v.sub {
			t.Errorf("Error in case %d, want %s, got %s", i, v.add, out)
		}
	}
}

func TestLength_Pixels(t *testing.T) {
	cases := []struct {
		in  Point
		dpi image.Point
		out image.Point
	}{
		{Point{1 * DIP, 2 * DIP}, image.Point{96, 96}, image.Point{1, 2}},
		{Point{1 * DIP, 2 * DIP}, image.Point{2 * 96, 3 * 96}, image.Point{2, 6}},
	}

	for i, v := range cases {
		DPI = v.dpi
		if out := v.in.Pixels(); out != v.out {
			t.Errorf("Error in case %d, want %s, got %s", i, v.out, out)
		}
	}
}

func TestRectangle(t *testing.T) {
	cases := []struct {
		x0, y0, x1, y1 Length
		min            Point
		width          Length
		height         Length
	}{
		{1 * DIP, 2 * DIP, 10 * DIP, 12 * DIP, Point{1 * DIP, 2 * DIP}, 9 * DIP, 10 * DIP},
		{1 * DIP, 12 * DIP, 10 * DIP, 2 * DIP, Point{1 * DIP, 2 * DIP}, 9 * DIP, 10 * DIP},
		{10 * DIP, 2 * DIP, 1 * DIP, 12 * DIP, Point{1 * DIP, 2 * DIP}, 9 * DIP, 10 * DIP},
		{10 * DIP, 12 * DIP, 1 * DIP, 2 * DIP, Point{1 * DIP, 2 * DIP}, 9 * DIP, 10 * DIP},
	}

	for i, v := range cases {
		out := Rect(v.x0, v.y0, v.x1, v.y1)
		if got := out.Dx(); got != v.width {
			t.Errorf("Error in case %d, want %s, got %s", i, got, v.width)
		}
		if got := out.Dy(); got != v.height {
			t.Errorf("Error in case %d, want %s, got %s", i, got, v.height)
		}
		expected := Point{v.width, v.height}
		if got := out.Size(); got != expected {
			t.Errorf("Error in case %d, want %s, got %s", i, got, expected)
		}
	}
}
