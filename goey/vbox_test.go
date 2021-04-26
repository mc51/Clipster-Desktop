package goey

import (
	"testing"

	"clipster/goey/base"
	"clipster/goey/mock"
)

func (w *vboxElement) Props() base.Widget {
	children := []base.Widget(nil)
	if len(w.children) != 0 {
		children = make([]base.Widget, 0, len(w.children))
		for _, v := range w.children {
			children = append(children, v.(Proper).Props())
		}
	}

	return &VBox{
		AlignMain:  w.alignMain,
		AlignCross: w.alignCross,
		Children:   children,
	}
}

func TestVBoxMount(t *testing.T) {
	buttons := []base.Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingMountWidgets(t,
		&VBox{},
		&VBox{Children: buttons, AlignMain: MainStart},
		&VBox{Children: buttons, AlignMain: MainCenter},
		&VBox{Children: buttons, AlignMain: MainEnd},
		&VBox{Children: buttons, AlignMain: SpaceAround},
		&VBox{Children: buttons, AlignMain: SpaceBetween},
		&VBox{Children: buttons, AlignMain: Homogeneous},
	)
}

func TestVBoxClose(t *testing.T) {
	buttons := []base.Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingCloseWidgets(t,
		&VBox{},
		&VBox{Children: buttons, AlignMain: MainStart},
	)
}

func TestVBoxUpdateProps(t *testing.T) {
	buttons := []base.Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingUpdateWidgets(t, []base.Widget{
		&VBox{AlignMain: MainStart},
		&VBox{Children: buttons, AlignMain: MainEnd, AlignCross: CrossStart},
	}, []base.Widget{
		&VBox{Children: buttons, AlignMain: MainEnd},
		&VBox{AlignMain: MainStart, AlignCross: CrossCenter},
	})
}

func TestVBoxLayout(t *testing.T) {
	children := mock.NewList(base.Size{10 * DIP, 26 * DIP}, base.Size{20 * DIP, 13 * DIP})

	cases := []struct {
		children    []base.Element
		alignMain   MainAxisAlign
		alignCross  CrossAxisAlign
		constraints base.Constraints
		size        base.Size
		bounds      []base.Rectangle
	}{
		{nil, MainStart, Stretch, base.TightWidth(40 * DIP), base.Size{40 * DIP, 0}, []base.Rectangle{}},
		{children, MainStart, Stretch, base.TightWidth(40 * DIP), base.Size{40 * DIP, 50 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 40*DIP, 26*DIP), base.Rect(0, 37*DIP, 40*DIP, 50*DIP),
		}},
		{children, MainCenter, Stretch, base.Tight(base.Size{40 * DIP, 150 * DIP}), base.Size{40 * DIP, 150 * DIP}, []base.Rectangle{
			base.Rect(0, 50*DIP, 40*DIP, 76*DIP), base.Rect(0, 87*DIP, 40*DIP, 100*DIP),
		}},
		{children, MainEnd, Stretch, base.Tight(base.Size{40 * DIP, 150 * DIP}), base.Size{40 * DIP, 150 * DIP}, []base.Rectangle{
			base.Rect(0, 100*DIP, 40*DIP, 126*DIP), base.Rect(0, 137*DIP, 40*DIP, 150*DIP),
		}},
		{children, SpaceBetween, Stretch, base.Tight(base.Size{40 * DIP, 150 * DIP}), base.Size{40 * DIP, 150 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 40*DIP, 26*DIP), base.Rect(0, 137*DIP, 40*DIP, 150*DIP),
		}},
		{children, MainStart, CrossStart, base.TightWidth(40 * DIP), base.Size{40 * DIP, 50 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 10*DIP, 26*DIP), base.Rect(0, 37*DIP, 20*DIP, 50*DIP),
		}},
		{children, MainStart, CrossCenter, base.TightWidth(40 * DIP), base.Size{40 * DIP, 50 * DIP}, []base.Rectangle{
			base.Rect(15*DIP, 0, 25*DIP, 26*DIP), base.Rect(10*DIP, 37*DIP, 30*DIP, 50*DIP),
		}},
		{children, MainStart, CrossEnd, base.TightWidth(40 * DIP), base.Size{40 * DIP, 50 * DIP}, []base.Rectangle{
			base.Rect(30*DIP, 0, 40*DIP, 26*DIP), base.Rect(20*DIP, 37*DIP, 40*DIP, 50*DIP),
		}},
		{children, Homogeneous, CrossStart, base.Tight(base.Size{40 * DIP, 150 * DIP}), base.Size{40 * DIP, 150 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 10*DIP, (2*75-11)*DIP/2), base.Rect(0, (2*75+11)*DIP/2, 20*DIP, 150*DIP),
		}},
		{children, Homogeneous, CrossCenter, base.Tight(base.Size{40 * DIP, 150 * DIP}), base.Size{40 * DIP, 150 * DIP}, []base.Rectangle{
			base.Rect(15*DIP, 0, 25*DIP, (2*75-11)*DIP/2), base.Rect(10*DIP, (2*75+11)*DIP/2, 30*DIP, 150*DIP),
		}},
		{children, Homogeneous, CrossEnd, base.Tight(base.Size{40 * DIP, 150 * DIP}), base.Size{40 * DIP, 150 * DIP}, []base.Rectangle{
			base.Rect(30*DIP, 0, 40*DIP, (2*75-11)*DIP/2), base.Rect(20*DIP, (2*75+11)*DIP/2, 40*DIP, 150*DIP),
		}},
		{children, Homogeneous, Stretch, base.Tight(base.Size{40 * DIP, 150 * DIP}), base.Size{40 * DIP, 150 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 40*DIP, (2*75-11)*DIP/2), base.Rect(0, (2*75+11)*DIP/2, 40*DIP, 150*DIP),
		}},
	}

	for i, v := range cases {
		in := vboxElement{
			children:     v.children,
			alignMain:    v.alignMain,
			alignCross:   v.alignCross,
			childrenInfo: make([]boxElementInfo, len(v.children)),
		}

		size := in.Layout(v.constraints)
		if size != v.size {
			t.Errorf("Incorrect size on case %d, got %s, want %s", i, size, v.size)
		}
		in.SetBounds(base.Rect(0, 0, size.Width, size.Height))
		for j, u := range v.bounds {
			if got := v.children[j].(*mock.Element).Bounds(); got != u {
				t.Errorf("Incorrect bounds case %d-%d, got %s, want %s", i, j, got, u)
			}
		}
	}
}

func TestVBoxMinIntrinsic(t *testing.T) {
	size := func(w, h base.Length) base.Size {
		return base.Size{w, h}
	}

	cases := []struct {
		children           []base.Element
		alignMain          MainAxisAlign
		alignCross         CrossAxisAlign
		minIntrinsicHeight base.Length
		minIntrinsicWidth  base.Length
	}{
		{nil, MainStart, Stretch, 0, 0},
		{mock.NewList(size(13*DIP, 13*DIP), size(13*DIP, 13*DIP)), MainStart, Stretch, 37 * DIP, 13 * DIP},
		{mock.NewList(size(13*DIP, 13*DIP), size(15*DIP, 13*DIP)), MainStart, Stretch, 37 * DIP, 15 * DIP},
		{mock.NewList(size(13*DIP, 26*DIP), size(11*DIP, 13*DIP)), MainStart, Stretch, 50 * DIP, 13 * DIP},
		{mock.NewList(size(13*DIP, 26*DIP), size(11*DIP, 13*DIP)), MainCenter, Stretch, 50 * DIP, 13 * DIP},
		{mock.NewList(size(13*DIP, 26*DIP), size(11*DIP, 13*DIP)), MainEnd, Stretch, 50 * DIP, 13 * DIP},
		{mock.NewList(size(13*DIP, 26*DIP), size(11*DIP, 13*DIP)), SpaceAround, Stretch, 72 * DIP, 13 * DIP},
		{mock.NewList(size(13*DIP, 26*DIP), size(11*DIP, 13*DIP)), SpaceBetween, Stretch, 50 * DIP, 13 * DIP},
		{mock.NewList(size(13*DIP, 26*DIP), size(11*DIP, 13*DIP)), Homogeneous, Stretch, (26*2 + 11) * DIP, 13 * DIP},
	}

	for i, v := range cases {
		in := vboxElement{
			children:   v.children,
			alignMain:  v.alignMain,
			alignCross: v.alignCross,
		}

		if value := in.MinIntrinsicHeight(base.Inf); value != v.minIntrinsicHeight {
			t.Errorf("Incorrect min intrinsic height on case %d, got %s, want %s", i, value, v.minIntrinsicHeight)
		}
		if value := in.MinIntrinsicWidth(base.Inf); value != v.minIntrinsicWidth {
			t.Errorf("Incorrect min intrinsic width on case %d, got %s, want %s", i, value, v.minIntrinsicWidth)
		}
	}
}
