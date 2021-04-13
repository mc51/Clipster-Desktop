package icons

import (
	"errors"
	"image"
	"image/draw"

	"golang.org/x/image/math/fixed"
)

var (
	ErrRuneNotAvailable = errors.New("rune not available")
)

// DrawImage returns a 32x32 image with the icon specified by the rune.
func DrawImage(index rune) (image.Image, error) {
	const Width = 32
	const Height = 32

	// Locate the index of this rune in the font file.
	ndx := assets.font.Index(index)
	if ndx == 0 {
		return nil, ErrRuneNotAvailable
	}

	// Measure geometry of rune to get placement, and then get the
	// masks for drawing.
	dr, _, _, _, _ := assets.face.Glyph(fixed.P(0, 0), index)
	dot := fixed.P(Width/2-dr.Dx()/2-dr.Min.X, Height/2+dr.Dy()/2-dr.Max.Y)
	dr, mask, maskp, _, _ := assets.face.Glyph(dot, index)

	// Draw the image.
	img := image.NewGray(image.Rect(0, 0, Width, Height))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Over)
	draw.DrawMask(img, dr, image.Black, image.Point{}, mask, maskp, draw.Over)
	return img, nil
}
