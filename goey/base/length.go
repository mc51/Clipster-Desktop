package base

import (
	"image"

	"golang.org/x/image/math/fixed"
)

var (
	// DPI contains the current DPI (dots per inch) of the monitor.
	// User code should not need to set this directly, as drivers will update
	// this variable as necessary.
	DPI image.Point
)

// Common lengths used when describing GUIs.
// Note that the DIP (device-independent pixel) is the natural unit for this
// package.  Because of limited precision, the PT listed here is somewhat smaller
// than its correct value.
const (
	DIP  Length = (1 << 6)         // Device-independent pixel (1/96 inch)
	PT   Length = ((96 << 6) / 72) // Point (1/72 inch)
	PC   Length = ((96 << 6) / 6)  // Pica (1/6 inch or 12 points)
	Inch Length = (96 << 6)        // Inch from the British imperial system of measurements
)

// Length is a distance measured in device-independent pixels.  There are nominally
// 96 DIPs per inch.  This definition corresponds with the definition of a
// pixel for both CSS and on Windows.
type Length fixed.Int26_6

// Clamp ensures that the length is between the minimum and maximum values
// specified.  Normally, min should be less than max.  If that is not the
// case, then the returned will preferentially respect min.
func (v Length) Clamp(min, max Length) Length {
	if v > max {
		v = max
	}
	if v < min {
		v = min
	}
	return v
}

// DIP returns a float64 with the length measured in device independent pixels.
func (v Length) DIP() float64 {
	return float64(v) / (1 << 6)
}

// Inch returns a float64 with the length measured in inches.
func (v Length) Inch() float64 {
	return float64(v) / (96 << 6)
}

// PT returns a float64 with the length measured in points.
func (v Length) PT() float64 {
	return float64(v) / ((96 << 6) / 72)
}

// PC returns a float64 with the length measured in picas.
func (v Length) PC() float64 {
	return float64(v) / ((96 << 6) / 6)
}

// PixelsX converts the distance measurement in DIPs to physical pixels, based
// on the current DPI settings for horizontal scaling.
func (v Length) PixelsX() int {
	return fixed.Int26_6(v.Scale(DPI.X, 96)).Round()
}

// PixelsY converts the distance measurement in DIPs to physical pixels, based
// on the current DPI settings for vertical scaling.
func (v Length) PixelsY() int {
	return fixed.Int26_6(v.Scale(DPI.Y, 96)).Round()
}

// Scale scales the distance by the ratio of num:den.
func (v Length) Scale(num, den int) Length {
	return Length(int64(v) * int64(num) / int64(den))
}

// String returns a human readable distance.
func (v Length) String() string {
	return fixed.Int26_6(v).String()
}

// FromPixelsX converts a distance measurement in physical pixels to DIPs, based on
// the current DPI settings for horizontal scaling.
func FromPixelsX(pixels int) Length {
	return Length(pixels<<6).Scale(96, DPI.X)
}

// FromPixelsY converts a distance measurement in physical pixels to DIPs, based on
// the current DPI settings for vertical scaling.
func FromPixelsY(pixels int) Length {
	return Length(pixels<<6).Scale(96, DPI.Y)
}

// A Point is an X, Y coordinate pair. The axes increase right and down.  This
// type is a close analogy to image.Pixel, except that the coordinate pair is
// represented by Length rather than int.
type Point struct {
	X, Y Length
}

// String returns a string representation of p like "(3,4)".
func (p Point) String() string {
	return "(" + p.X.String() + "," + p.Y.String() + ")"
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Pixels returns the vector with the X and Y coordinates measured in pixels.
func (p Point) Pixels() image.Point {
	return image.Point{p.X.PixelsX(), p.Y.PixelsY()}
}

// A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
//
// This type is a close analogy to image.Rectangle, except that the coordinates
// are represented by Length rather than int.
type Rectangle struct {
	Min, Max Point
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Add returns the rectangle r translated by p.
func (r Rectangle) Add(p Point) Rectangle {
	return Rectangle{
		Point{r.Min.X + p.X, r.Min.Y + p.Y},
		Point{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Dx returns r's width.
func (r Rectangle) Dx() Length {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle) Dy() Length {
	return r.Max.Y - r.Min.Y
}

// Pixels returns the rectangle with the X and Y coordinates measured in pixels.
func (r Rectangle) Pixels() image.Rectangle {
	return image.Rectangle{r.Min.Pixels(), r.Max.Pixels()}
}

// Size returns r's width and height.
func (r Rectangle) Size() Point {
	return Point{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Rect is shorthand for Rectangle{Point(x0, y0), Point(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func Rect(x0, y0, x1, y1 Length) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}

	if y0 > y1 {
		y0, y1 = y1, y0
	}

	return Rectangle{Point{x0, y0}, Point{x1, y1}}
}
