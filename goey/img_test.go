package goey

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"guitest/goey/base"
)

func drawVerticalRGB(img draw.Image) {
	colors := [3]color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
	draw.Draw(img, image.Rect(0, 0, dx/3, dy), image.NewUniform(colors[0]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(dx/3, 0, dx*2/3, dy), image.NewUniform(colors[1]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(dx*2/3, 0, dx, dy), image.NewUniform(colors[2]), image.Point{}, draw.Src)
}

func drawHorizontalRGB(img draw.Image) {
	colors := [3]color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
	draw.Draw(img, image.Rect(0, 0, dx, dy/3), image.NewUniform(colors[0]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(0, dy/3, dx, dy*2/3), image.NewUniform(colors[1]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(0, dy*2/3, dx, dy), image.NewUniform(colors[2]), image.Point{}, draw.Src)
}

func drawVerticalGradient(img *image.Gray) {
	h := img.Rect.Dy()
	for i := img.Rect.Min.Y; i < img.Rect.Max.Y; i++ {
		clr := (i - img.Rect.Min.Y) * 255 / h
		for j := img.Rect.Min.X; j < img.Rect.Max.X; j++ {
			img.Pix[img.PixOffset(j, i)] = uint8(clr)
		}
	}
}

func TestImgMount(t *testing.T) {
	bounds := image.Rect(0, 0, 92, 92)
	images := []draw.Image{image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds),
		image.NewGray(bounds), image.NewGray(bounds)}
	draw.Draw(images[0], bounds, image.NewUniform(color.RGBA{255, 255, 0, 255}), image.Point{}, draw.Src)
	draw.Draw(images[1], bounds, image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(images[2], bounds, image.NewUniform(color.RGBA{0, 255, 255, 255}), image.Point{}, draw.Src)
	drawVerticalRGB(images[3])
	drawHorizontalRGB(images[4])
	draw.Draw(images[5], bounds, image.NewUniform(color.RGBA{128, 128, 128, 255}), image.Point{}, draw.Src)
	drawVerticalGradient(images[6].(*image.Gray))

	testingMountWidgets(t,
		&Align{Child: &Img{Image: images[0], Width: 100 * DIP, Height: 10 * DIP}},
		&Align{Child: &Img{Image: images[1], Width: 20 * DIP}},
		&Align{Child: &Img{Image: images[2], Height: 30 * DIP}},
		&Align{Child: &Img{Image: images[3]}},
		&Align{Child: &Img{Image: images[4]}},
		&Align{Child: &Img{Image: images[5]}},
		&Align{Child: &Img{Image: images[6]}},
	)
}

func TestImgClose(t *testing.T) {
	bounds := image.Rect(0, 0, 92, 92)
	images := []*image.RGBA{image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds)}
	draw.Draw(images[0], bounds, image.NewUniform(color.RGBA{255, 255, 0, 255}), image.Point{}, draw.Src)
	draw.Draw(images[1], bounds, image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(images[2], bounds, image.NewUniform(color.RGBA{0, 255, 255, 255}), image.Point{}, draw.Src)

	testingCloseWidgets(t,
		&Img{Image: images[0], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[1], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[2]},
	)
}

func TestImgMinIntrinsicHeight(t *testing.T) {
	cases := []struct {
		width   base.Length
		height  base.Length
		atWidth base.Length
		out     base.Length
	}{
		{10 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{10 * DIP, 10 * DIP, 20 * DIP, 10 * DIP},
		{10 * DIP, 10 * DIP, base.Inf, 10 * DIP},
		{10 * DIP, 10 * DIP, 5 * DIP, 5 * DIP},
		{10 * DIP, 20 * DIP, 10 * DIP, 20 * DIP},
		{10 * DIP, 20 * DIP, 20 * DIP, 20 * DIP},
		{10 * DIP, 20 * DIP, base.Inf, 20 * DIP},
		{10 * DIP, 20 * DIP, 5 * DIP, 10 * DIP},
		{20 * DIP, 10 * DIP, 10 * DIP, 5 * DIP},
		{20 * DIP, 10 * DIP, 20 * DIP, 10 * DIP},
		{20 * DIP, 10 * DIP, base.Inf, 10 * DIP},
		{20 * DIP, 10 * DIP, 5 * DIP, 5 * DIP / 2},
	}

	for i, v := range cases {
		elem := imgElement{
			width:  v.width,
			height: v.height,
		}

		if out := elem.MinIntrinsicHeight(v.atWidth); out != v.out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
	}
}

func TestImgMinIntrinsicWidth(t *testing.T) {
	cases := []struct {
		width    base.Length
		height   base.Length
		atHeight base.Length
		out      base.Length
	}{
		{10 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{10 * DIP, 10 * DIP, 20 * DIP, 10 * DIP},
		{10 * DIP, 10 * DIP, base.Inf, 10 * DIP},
		{10 * DIP, 10 * DIP, 5 * DIP, 5 * DIP},
		{20 * DIP, 10 * DIP, 10 * DIP, 20 * DIP},
		{20 * DIP, 10 * DIP, 20 * DIP, 20 * DIP},
		{20 * DIP, 10 * DIP, base.Inf, 20 * DIP},
		{20 * DIP, 10 * DIP, 5 * DIP, 10 * DIP},
		{10 * DIP, 20 * DIP, 10 * DIP, 5 * DIP},
		{10 * DIP, 20 * DIP, 20 * DIP, 10 * DIP},
		{10 * DIP, 20 * DIP, base.Inf, 10 * DIP},
		{10 * DIP, 20 * DIP, 5 * DIP, 5 * DIP / 2},
	}

	for i, v := range cases {
		elem := imgElement{
			width:  v.width,
			height: v.height,
		}

		if out := elem.MinIntrinsicWidth(v.atHeight); out != v.out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
	}
}

func TestImgUpdate(t *testing.T) {
	bounds := image.Rect(0, 0, 92, 92)
	images := []*image.RGBA{image.NewRGBA(bounds), image.NewRGBA(bounds), image.NewRGBA(bounds)}
	draw.Draw(images[0], bounds, image.NewUniform(color.RGBA{255, 255, 0, 255}), image.Point{}, draw.Src)
	draw.Draw(images[1], bounds, image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(images[2], bounds, image.NewUniform(color.RGBA{0, 255, 255, 255}), image.Point{}, draw.Src)

	testingUpdateWidgets(t, []base.Widget{
		&Img{Image: images[0], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[1], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[2], Width: 100 * DIP, Height: 10 * DIP},
	}, []base.Widget{
		&Img{Image: images[2], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[1], Width: 100 * DIP, Height: 10 * DIP},
		&Img{Image: images[0], Width: 100 * DIP, Height: 10 * DIP},
	})
}

func TestImgUpdateDimensions(t *testing.T) {
	img1 := image.RGBA{Rect: image.Rect(0, 0, 92, 92)}

	cases := []struct {
		width  base.Length
		height base.Length
		img    image.Image
		out    base.Size
	}{
		{10 * DIP, 15 * DIP, &img1, base.Size{10 * DIP, 15 * DIP}},
		{0, 0, &img1, base.Size{1 * Inch, 1 * Inch}},
		{2 * Inch, 0, &img1, base.Size{2 * Inch, 2 * Inch}},
		{0, 2 * Inch, &img1, base.Size{2 * Inch, 2 * Inch}},
	}

	for i, v := range cases {
		widget := Img{
			Width:  v.width,
			Height: v.height,
			Image:  v.img,
		}

		widget.UpdateDimensions()
		if widget.Height != v.out.Height {
			t.Errorf("Case %d:  Failed to update height, got %v, want %v", i, widget.Height, v.out.Height)
		}
		if widget.Width != v.out.Width {
			t.Errorf("Case %d:  Failed to update width, got %v, want %v", i, widget.Width, v.out.Width)
		}
	}
}
