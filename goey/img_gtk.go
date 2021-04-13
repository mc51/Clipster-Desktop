// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"image"
	"image/draw"

	"guitest/goey/base"
	"guitest/goey/internal/gtk"
)

type imgElement struct {
	Control

	pix    []uint8
	width  base.Length
	height base.Length
}

func imageToRGBA(prop image.Image) *image.RGBA {
	// Use existing image if possible
	if img, ok := prop.(*image.RGBA); ok {
		return &image.RGBA{
			Pix:    append([]uint8(nil), img.Pix...),
			Stride: img.Stride,
			Rect:   img.Rect,
		}
	}

	// Create a new image in RGBA format
	bounds := prop.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, prop, bounds.Min, draw.Src)
	return img
}

func (w *Img) mount(parent base.Control) (base.Element, error) {
	img := imageToRGBA(w.Image)

	handle := gtk.MountImage(parent.Handle, &img.Pix[0], img.Rect.Dx(), img.Rect.Dy(), img.Stride)

	retval := &imgElement{
		Control: Control{handle},
		pix:     img.Pix,
		width:   w.Width,
		height:  w.Height,
	}
	gtk.RegisterWidget(handle, retval)

	return retval, nil
}

func (w *imgElement) Props() base.Widget {
	return &Img{
		Image:  w.PropsImage(),
		Width:  w.width,
		Height: w.height,
	}
}

func (w *imgElement) PropsImage() image.Image {
	// Assuming 8-bits per pixel.

	if gtk.ImageColorSpace(w.handle) != 0 {
		return nil
	}

	width := int(gtk.ImageImageWidth(w.handle))
	height := int(gtk.ImageImageHeight(w.handle))
	stride := int(gtk.ImageImageStride(w.handle))

	// If there is an alpha channel, we can use the data directly.
	if gtk.ImageHasAlpha(w.handle) {
		return &image.RGBA{
			Pix:    gtk.ImageImageData(w.handle),
			Stride: stride,
			Rect:   image.Rect(0, 0, width, height),
		}
	}

	// Need munge the pixel data into the correct format.
	// Convert RGB to RGBA.
	pix := gtk.ImageImageData(w.handle)
	newpix := make([]byte, height*width*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newpix[y*width*4+x*4+0] = pix[y*stride+x*3+0]
			newpix[y*width*4+x*4+1] = pix[y*stride+x*3+1]
			newpix[y*width*4+x*4+2] = pix[y*stride+x*3+2]
			newpix[y*width*4+x*4+3] = 0xff
		}
	}

	// Note:  stride of the new image data does not match data returned
	// from Pixbuf.
	return &image.RGBA{
		Pix:    newpix,
		Stride: width * 4,
		Rect:   image.Rect(0, 0, width, height),
	}
}

func (w *imgElement) updateProps(data *Img) error {
	w.width, w.height = data.Width, data.Height

	// Create the bitmap
	img := imageToRGBA(data.Image)
	gtk.ImageUpdate(w.handle, &img.Pix[0], img.Rect.Dx(), img.Rect.Dy(), img.Stride)
	w.pix = img.Pix

	return nil
}
