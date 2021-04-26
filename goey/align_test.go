package goey

import (
	"errors"
	"testing"

	"clipster/goey/base"
	"clipster/goey/mock"
)

func (w *alignElement) Props() base.Widget {
	child := base.Widget(nil)
	if w.child != nil {
		child = w.child.(Proper).Props()
	}

	return &Align{
		HAlign:       w.hAlign,
		VAlign:       w.vAlign,
		WidthFactor:  w.widthFactor,
		HeightFactor: w.heightFactor,
		Child:        child,
	}
}

func TestAlignMount(t *testing.T) {
	// These should all be able to mount without error.
	testingMountWidgets(t,
		&Align{Child: &Button{Text: "A"}},
		&Align{HAlign: AlignStart, Child: &Button{Text: "B"}},
		&Align{HAlign: AlignEnd, Child: &Button{Text: "C"}},
		&Align{HAlign: AlignCenter, Child: &Button{Text: "C"}},
		&Align{HeightFactor: 2, WidthFactor: 2.5, Child: &Button{Text: "C"}},
		&Align{},
	)

	// These should mount with an error.
	err := errors.New("Mock error 1")
	testingMountWidgetsFail(t, err,
		&Align{Child: &mock.Widget{Err: err}},
	)
}

func TestAlignClose(t *testing.T) {
	testingCloseWidgets(t,
		&Align{Child: &Button{Text: "A"}},
		&Align{},
	)
}

func TestAlignUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Align{Child: &Button{Text: "A"}},
		&Align{HAlign: AlignStart, Child: &Button{Text: "B"}},
		&Align{HAlign: AlignEnd, Child: &Button{Text: "C"}},
		&Align{HAlign: AlignCenter, Child: &Button{Text: "C"}},
		&Align{HeightFactor: 2, WidthFactor: 2.5, Child: &Button{Text: "C"}},
		&Align{},
	}, []base.Widget{
		&Align{Child: &Button{Text: "AB"}},
		&Align{HAlign: AlignCenter, Child: &Button{Text: "BC"}},
		&Align{HAlign: AlignStart, Child: &Button{Text: "CD"}},
		&Align{HAlign: AlignEnd, Child: &Button{Text: "CE"}},
		&Align{HeightFactor: 4, WidthFactor: 3},
		&Align{Child: &Label{Text: "CF"}},
	})
}

func TestAlignLayout(t *testing.T) {
	size1 := base.Size{10 * DIP, 20 * DIP}
	size2 := base.Size{35 * DIP, 30 * DIP}
	sizeZ := base.Size{}

	cases := []struct {
		in           base.Size
		widthFactor  float64
		heightFactor float64
		bc           base.Constraints
		out          base.Size
	}{
		{size1, 0, 0, base.Loose(size1), size1},
		{size2, 0, 0, base.Loose(size2), size2},
		{size1, 0, 0, base.Loose(size2), size1},
		{sizeZ, 0, 0, base.Loose(size1), size1},
		{sizeZ, 0, 0, base.Loose(size2), size2},
		{size1, 2, 2, base.Loose(size1), size1},
		{size2, 2, 2, base.Loose(size2), size2},
		{size1, 2, 2, base.Loose(size2), base.Size{20 * DIP, 30 * DIP}},
		{size1, 0, 0, base.Tight(size1), size1},
		{size2, 0, 0, base.Tight(size2), size2},
		{sizeZ, 0, 0, base.Tight(size1), size1},
		{size1, 0, 0, base.Expand().Loosen(), size1},
		{size2, 0, 0, base.Expand().Loosen(), size2},
		{sizeZ, 0, 0, base.Expand().Loosen(), sizeZ},
		{size1, 1, 1, base.Expand().Loosen(), size1},
		{size2, 1, 1, base.Expand().Loosen(), size2},
		{size1, 2, 3, base.Expand().Loosen(), base.Size{20 * DIP, 60 * DIP}},
		{size2, 4, 5, base.Expand().Loosen(), base.Size{4 * 35 * DIP, 150 * DIP}},
	}

	for i, v := range cases {
		elem := alignElement{
			child:        mock.NewIfNotZero(v.in),
			widthFactor:  v.widthFactor,
			heightFactor: v.heightFactor,
		}

		if out := elem.Layout(v.bc); out != v.out {
			t.Errorf("Case %d: Returned size does not match, got %v, want %v", i, out, v.out)
		}
	}
}

func TestAlignMinIntrinsicSize(t *testing.T) {
	size1 := base.Size{10 * DIP, 20 * DIP}
	size2 := base.Size{35 * DIP, 30 * DIP}
	sizeZ := base.Size{}

	cases := []struct {
		in           base.Size
		widthFactor  float64
		heightFactor float64
		out          base.Size
	}{
		{size1, 0, 0, size1},
		{size2, 0, 0, size2},
		{sizeZ, 0, 0, sizeZ},
		{size1, 1, 1, size1},
		{size2, 1, 1, size2},
		{sizeZ, 1, 1, sizeZ},
		{size1, 2, 3, base.Size{20 * DIP, 60 * DIP}},
	}

	for i, v := range cases {
		elem := alignElement{
			child:        mock.NewIfNotZero(v.in),
			widthFactor:  v.widthFactor,
			heightFactor: v.heightFactor,
		}

		if out := elem.MinIntrinsicWidth(base.Inf); out != v.out.Width {
			t.Errorf("Case %d: Returned min intrinsic width does not match, got %v, want %v", i, out, v.out.Width)
		}
		if out := elem.MinIntrinsicHeight(base.Inf); out != v.out.Height {
			t.Errorf("Case %d: Returned min intrinsic width does not match, got %v, want %v", i, out, v.out.Height)
		}
	}
}
