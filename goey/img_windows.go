package goey

import (
	"errors"
	"image"
	"image/draw"
	"unsafe"

	"guitest/goey/base"
	win2 "guitest/goey/internal/syscall"
	"github.com/lxn/win"
)

func imageToIcon(prop image.Image) (win.HICON, []uint8, error) {
	// Create a mask for the icon.
	// Currently, we are using a straight white mask, but perhaps this
	// should be a copy of the alpha channel if the source image is
	// RGBA.
	bounds := prop.Bounds()
	imgMask := image.NewGray(prop.Bounds())
	draw.Draw(imgMask, bounds, image.White, image.Point{}, draw.Src)
	hmask, _, err := imageToBitmap(imgMask)
	if err != nil {
		return 0, nil, err
	}

	// Convert the image to a bitmap.
	hbitmap, buffer, err := imageToBitmap(prop)
	if err != nil {
		return 0, nil, err
	}

	// Create the icon
	iconinfo := win.ICONINFO{
		FIcon:    win.TRUE,
		HbmMask:  hmask,
		HbmColor: hbitmap,
	}
	hicon := win.CreateIconIndirect(&iconinfo)
	if hicon == 0 {
		panic("Error in CreateIconIndirect")
	}
	return hicon, buffer, nil
}

func imageToBitmapRGBA(img *image.RGBA, pix []byte) (win.HBITMAP, error) {
	// Need to convert RGBA to BGRA.
	for i := 0; i < len(pix); i += 4 {
		// swap the red and green bytes.
		pix[i+0], pix[i+2] = pix[i+2], pix[i+0]
	}

	// The following call also works with 4 channels of 8-bits on a Windows
	// machine, but fails on Wine.  Would like it to work on both to ease
	// CI.
	hbitmap := win.CreateBitmap(int32(img.Rect.Dx()), int32(img.Rect.Dy()), 1, 32, unsafe.Pointer(&pix[0]))
	if hbitmap == 0 {
		return 0, errors.New("call to CreateBitmap failed")
	}
	return hbitmap, nil
}

func imageToBitmap(prop image.Image) (win.HBITMAP, []uint8, error) {
	if img, ok := prop.(*image.RGBA); ok {
		// Create a copy of the backing for the pixel data
		buffer := append([]uint8(nil), img.Pix...)
		// Create the bitmap.
		hbitmap, err := imageToBitmapRGBA(img, buffer)
		return hbitmap, buffer, err
	}

	// Create a new image in RGBA format
	bounds := prop.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, prop, bounds.Min, draw.Src)
	// Create the bitmap
	hbitmap, err := imageToBitmapRGBA(img, img.Pix)
	return hbitmap, img.Pix, err
}

func bitmapToImage(hdc win.HDC, hbitmap win.HBITMAP) image.Image {
	bmi := win.BITMAPINFO{}
	bmi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bmi))
	win.GetDIBits(hdc, hbitmap, 0, 0, nil, &bmi, 0)
	if bmi.BmiHeader.BiPlanes == 1 && bmi.BmiHeader.BiBitCount == 32 && bmi.BmiHeader.BiCompression == win.BI_BITFIELDS {
		// Get the pixel data
		buffer := make([]byte, bmi.BmiHeader.BiSizeImage)
		win.GetDIBits(hdc, hbitmap, 0, uint32(bmi.BmiHeader.BiHeight), &buffer[0], &bmi, 0)

		// Need to convert BGR to RGB
		for i := 0; i < len(buffer); i += 4 {
			buffer[i+0], buffer[i+2] = buffer[i+2], buffer[i+0]
		}
		// In GDI, all bitmaps are bottom up.  We need to reorder the rows
		// before the data can be used for a PNG.
		// TODO:  Combine this pass with the previous.
		stride := int(bmi.BmiHeader.BiWidth) * 4
		for y := 0; y < int(bmi.BmiHeader.BiHeight/2); y++ {
			y2 := int(bmi.BmiHeader.BiHeight) - y - 1
			for x := 0; x < stride; x++ {
				// The stride is always the same as the width?
				buffer[y*stride+x], buffer[y2*stride+x] = buffer[y2*stride+x], buffer[y*stride+x]
			}
		}
		return &image.RGBA{
			Pix:    buffer,
			Stride: int(bmi.BmiHeader.BiWidth * 4),
			Rect:   image.Rect(0, 0, int(bmi.BmiHeader.BiWidth), int(bmi.BmiHeader.BiHeight)),
		}
	}

	return nil
}

func (w *Img) mount(parent base.Control) (base.Element, error) {
	// Create the bitmap
	hbitmap, buffer, err := imageToBitmap(w.Image)
	if err != nil {
		return nil, err
	}

	// Create the control
	const STYLE = win.WS_CHILD | win.WS_VISIBLE | win.SS_BITMAP | win.SS_LEFT
	hwnd, _, err := createControlWindow(0, &staticClassName[0], "", STYLE, parent.HWnd)
	if err != nil {
		return nil, err
	}
	win.SendMessage(hwnd, win2.STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	retval := &imgElement{
		Control:   Control{hwnd},
		imageData: buffer,
		width:     w.Width,
		height:    w.Height,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type imgElement struct {
	Control
	imageData []uint8
	width     base.Length
	height    base.Length
}

func (w *imgElement) Props() base.Widget {
	// Need to recreate the image from the HBITMAP
	hbitmap := win.HBITMAP(win.SendMessage(w.hWnd, win2.STM_GETIMAGE, 0 /*IMAGE_BITMAP*/, 0))
	if hbitmap == 0 {
		return &Img{
			Width:  w.width,
			Height: w.height,
		}
	}

	hdc := win.GetDC(w.hWnd)
	img := bitmapToImage(hdc, hbitmap)
	win.ReleaseDC(w.hWnd, hdc)

	return &Img{
		Image:  img,
		Width:  w.width,
		Height: w.height,
	}
}

func (w *imgElement) SetBounds(bounds base.Rectangle) {
	w.Control.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *imgElement) updateProps(data *Img) error {
	w.width, w.height = data.Width, data.Height

	// Create the bitmap
	hbitmap, buffer, err := imageToBitmap(data.Image)
	if err != nil {
		return err
	}
	w.imageData = buffer
	win.SendMessage(w.hWnd, win2.STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	return nil
}
