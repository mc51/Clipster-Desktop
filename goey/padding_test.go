package goey

import (
	"errors"
	"fmt"
	"testing"

	"guitest/goey/base"
	"guitest/goey/mock"
)

func (w *paddingElement) Props() base.Widget {
	child := base.Widget(nil)
	if w.child != nil {
		child = w.child.(Proper).Props()
	}

	return &Padding{
		Insets: w.insets,
		Child:  child,
	}
}

func ExampleInsets_String() {
	p := Insets{
		Top:    2 * base.DIP,
		Right:  1 * base.DIP,
		Bottom: 3 * base.DIP / 2,
		Left:   1 * base.DIP,
	}

	fmt.Println("Insets:", p)

	// Output:
	// Insets: {2:00 1:00 1:32 1:00}
}

func TestPaddingMount(t *testing.T) {
	// These should all be able to mount without error.
	testingMountWidgets(t,
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
		&Padding{},
	)

	// These should mount with an error.
	err := errors.New("Mock error 1")
	testingMountWidgetsFail(t, err,
		&Padding{Child: &mock.Widget{Err: err}},
	)
	testingMountWidgetsFail(t, err,
		&Padding{Insets: DefaultInsets(), Child: &mock.Widget{Err: err}},
	)
}

func TestPaddingClose(t *testing.T) {
	testingCloseWidgets(t,
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
		&Padding{},
	)
}

func TestPaddingUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
		&Padding{},
	}, []base.Widget{
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "AB"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "BC"}},
		&Padding{},
		&Padding{Child: &Button{Text: "CD"}},
	})
}

func TestPaddingMinIntrinsicSize(t *testing.T) {
	size1 := base.Size{10 * DIP, 20 * DIP}
	size2 := base.Size{15 * DIP, 25 * DIP}
	sizeZ := base.Size{}
	insets := Insets{1 * DIP, 2 * DIP, 3 * DIP, 4 * DIP}

	cases := []struct {
		mockSize base.Size
		insets   Insets
		out      base.Size
	}{
		{size1, Insets{}, size1},
		{size2, Insets{}, size2},
		{sizeZ, Insets{}, sizeZ},
		{size1, insets, base.Size{16 * DIP, 24 * DIP}},
		{size2, insets, base.Size{21 * DIP, 29 * DIP}},
		{sizeZ, insets, base.Size{6 * DIP, 4 * DIP}},
	}

	for i, v := range cases {
		elem := paddingElement{
			child:  mock.NewIfNotZero(v.mockSize),
			insets: v.insets,
		}

		if out := elem.MinIntrinsicWidth(base.Inf); out != v.out.Width {
			t.Errorf("Case %d: Returned min intrinsic width does not match, got %v, want %v", i, out, v.out.Width)
		}
		if out := elem.MinIntrinsicHeight(base.Inf); out != v.out.Height {
			t.Errorf("Case %d: Returned min intrinsic width does not match, got %v, want %v", i, out, v.out.Height)
		}
	}
}
