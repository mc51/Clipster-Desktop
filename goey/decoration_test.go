package goey

import (
	"errors"
	"image/color"
	"testing"

	"clipster/goey/base"
	"clipster/goey/mock"
)

var (
	black = color.RGBA{0, 0, 0, 0xff}
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
	red   = color.RGBA{0xcc, 0xaa, 0x88, 0xff}
)

func (w *decorationElement) Props() base.Widget {
	widget := w.props()
	if w.child != nil {
		widget.Child = w.child.(Proper).Props()
	}

	return widget
}

func TestDecorationMount(t *testing.T) {
	child := &mock.Widget{Size: base.Size{15 * base.DIP, 15 * base.DIP}}

	// These should all be able to mount without error.
	testingMountWidgets(t,
		&Decoration{Child: &Button{Text: "A"}, Insets: DefaultInsets()},
		&Decoration{Child: &Label{Text: "A"}, Fill: red, Insets: DefaultInsets()},
		&Decoration{Child: child},
		&Decoration{Child: child, Stroke: black},
		&Decoration{Child: child, Stroke: white},
		&Decoration{Child: child, Stroke: red},
		&Decoration{Child: child, Fill: black, Stroke: white, Radius: 8 * DIP, Insets: DefaultInsets()},
		&Decoration{Child: child, Fill: white, Stroke: black},
		&Decoration{Child: child, Fill: red},
	)

	// These should mount with an error.
	err := errors.New("Mock error 1")
	testingMountWidgetsFail(t, err,
		&Decoration{Child: &mock.Widget{Err: err}},
	)
	testingMountWidgetsFail(t, err,
		&Decoration{Insets: DefaultInsets(), Child: &mock.Widget{Err: err}},
	)
}

func TestDecorationClose(t *testing.T) {
	testingCloseWidgets(t,
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
	)
}

func TestDecorationUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
		&Decoration{Fill: white, Stroke: black},
		&Decoration{Fill: red},
	}, []base.Widget{
		&Decoration{},
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: black},
		&Decoration{Fill: white},
	})
}

func TestDecorationMinIntrinsicSize(t *testing.T) {
	size1 := base.Size{10 * DIP, 20 * DIP}
	size2 := base.Size{15 * DIP, 25 * DIP}
	sizeZ := base.Size{0, 0}
	sizeMin := base.Size{base.FromPixelsX(2), base.FromPixelsY(2)}
	insets := Insets{1 * DIP, 2 * DIP, 3 * DIP, 4 * DIP}

	cases := []struct {
		mockSize base.Size
		insets   Insets
		out      base.Size
	}{
		{size1, Insets{}, size1},
		{size2, Insets{}, size2},
		{sizeZ, Insets{}, sizeMin},
		{size1, insets, base.Size{16 * DIP, 24 * DIP}},
		{size2, insets, base.Size{21 * DIP, 29 * DIP}},
		{sizeZ, insets, base.Size{6 * DIP, 4 * DIP}},
	}

	for i, v := range cases {
		elem := decorationElement{
			child:  mock.NewIfNotZero(v.mockSize),
			insets: v.insets,
		}

		if out := elem.MinIntrinsicWidth(base.Inf); out != v.out.Width {
			t.Errorf("Case %d: Returned min intrinsic width does not match, got %v, want %v", i, out, v.out.Width)
		}
		if out := elem.MinIntrinsicHeight(base.Inf); out != v.out.Height {
			t.Errorf("Case %d: Returned min intrinsic height does not match, got %v, want %v", i, out, v.out.Height)
		}
	}
}
