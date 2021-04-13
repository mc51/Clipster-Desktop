package base

// Size represents the size of a rectangular element.
type Size struct {
	Width, Height Length
}

// FromPixels converts the pixels into lengths based on the current DPI, and
// return the size.
func FromPixels(x, y int) Size {
	return Size{FromPixelsX(x), FromPixelsY(y)}
}

// IsZero returns true if the size is the zero value.
func (s *Size) IsZero() bool {
	return s.Width == 0 && s.Height == 0
}

// String returns a string representation of the size.
func (s *Size) String() string {
	return "(" + s.Width.String() + "x" + s.Height.String() + ")"
}
