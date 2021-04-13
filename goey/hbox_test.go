package goey

import (
	"testing"

	"guitest/goey/base"
	"guitest/goey/mock"
)

func (w *hboxElement) Props() base.Widget {
	children := []base.Widget(nil)
	if len(w.children) != 0 {
		children = make([]base.Widget, 0, len(w.children))
		for _, v := range w.children {
			children = append(children, v.(Proper).Props())
		}
	}

	return &HBox{
		AlignMain:  w.alignMain,
		AlignCross: w.alignCross,
		Children:   children,
	}
}

func TestHBoxMount(t *testing.T) {
	buttons := []base.Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingMountWidgets(t,
		&HBox{},
		&HBox{Children: buttons, AlignMain: MainStart},
		&HBox{Children: buttons, AlignMain: MainCenter},
		&HBox{Children: buttons, AlignMain: MainEnd},
		&HBox{Children: buttons, AlignMain: SpaceAround},
		&HBox{Children: buttons, AlignMain: SpaceBetween},
		&HBox{Children: buttons, AlignMain: Homogeneous},
	)
}

func TestHBoxClose(t *testing.T) {
	buttons := []base.Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingCloseWidgets(t,
		&HBox{},
		&HBox{Children: buttons, AlignMain: MainStart},
	)
}

func TestHBoxUpdateProps(t *testing.T) {
	buttons := []base.Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}
	widgets1 := []base.Widget{
		&Button{Text: "A"},
		&Label{Text: "B"},
		&Button{Text: "C"},
	}
	widgets2 := []base.Widget{
		&Label{Text: "AA"},
		&Button{Text: "BA"},
		&Label{Text: "CA"},
		&Button{Text: "DA"},
	}

	testingUpdateWidgets(t, []base.Widget{
		&HBox{AlignMain: MainStart},
		&HBox{Children: buttons, AlignMain: MainEnd, AlignCross: CrossStart},
		&HBox{Children: widgets1},
		&HBox{Children: widgets2},
	}, []base.Widget{
		&HBox{Children: buttons, AlignMain: MainEnd},
		&HBox{AlignMain: MainStart, AlignCross: CrossCenter},
		&HBox{Children: widgets2},
		&HBox{Children: widgets1},
	})
}

func TestHBoxLayout(t *testing.T) {
	children := []base.Element{
		mock.New(base.Size{26 * DIP, 13 * DIP}), mock.New(base.Size{13 * DIP, 11 * DIP})}

	cases := []struct {
		children    []base.Element
		alignMain   MainAxisAlign
		alignCross  CrossAxisAlign
		constraints base.Constraints
		size        base.Size
		bounds      []base.Rectangle
	}{
		{nil, MainStart, Stretch, base.TightHeight(40 * DIP), base.Size{0, 40 * DIP}, []base.Rectangle{}},
		{children, MainStart, Stretch, base.TightHeight(40 * DIP), base.Size{50 * DIP, 40 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 26*DIP, 40*DIP), base.Rect(37*DIP, 0, 50*DIP, 40*DIP),
		}},
		{children, MainStart, Stretch, base.Tight(base.Size{150 * DIP, 40 * DIP}), base.Size{150 * DIP, 40 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 26*DIP, 40*DIP), base.Rect(37*DIP, 0, 50*DIP, 40*DIP),
		}},
		{children, MainEnd, Stretch, base.Tight(base.Size{150 * DIP, 40 * DIP}), base.Size{150 * DIP, 40 * DIP}, []base.Rectangle{
			base.Rect(100*DIP, 0, 126*DIP, 40*DIP), base.Rect(137*DIP, 0, 150*DIP, 40*DIP),
		}},
		{children, SpaceBetween, Stretch, base.Tight(base.Size{150 * DIP, 40 * DIP}), base.Size{150 * DIP, 40 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 26*DIP, 40*DIP), base.Rect(137*DIP, 0, 150*DIP, 40*DIP),
		}},
	}

	for i, v := range cases {
		in := hboxElement{
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

func TestHBoxMinIntrinsic(t *testing.T) {
	size := func(w, h base.Length) base.Size {
		return base.Size{w, h}
	}

	cases := []struct {
		children           []base.Element
		alignMain          MainAxisAlign
		alignCross         CrossAxisAlign
		minIntrinsicWidth  base.Length
		minIntrinsicHeight base.Length
	}{
		{nil, MainStart, Stretch, 0, 0},
		{mock.NewList(size(13*DIP, 13*DIP), size(13*DIP, 13*DIP)), MainStart, Stretch, 37 * DIP, 13 * DIP},
		{mock.NewList(size(13*DIP, 13*DIP), size(13*DIP, 15*DIP)), MainStart, Stretch, 37 * DIP, 15 * DIP},
		{mock.NewList(size(26*DIP, 13*DIP), size(13*DIP, 11*DIP)), MainStart, Stretch, 50 * DIP, 13 * DIP},
		{mock.NewList(size(26*DIP, 13*DIP), size(13*DIP, 11*DIP)), MainCenter, Stretch, 50 * DIP, 13 * DIP},
		{mock.NewList(size(26*DIP, 13*DIP), size(13*DIP, 11*DIP)), MainEnd, Stretch, 50 * DIP, 13 * DIP},
		{mock.NewList(size(26*DIP, 13*DIP), size(13*DIP, 11*DIP)), SpaceAround, Stretch, 72 * DIP, 13 * DIP},
		{mock.NewList(size(26*DIP, 13*DIP), size(13*DIP, 11*DIP)), SpaceBetween, Stretch, 50 * DIP, 13 * DIP},
		{mock.NewList(size(26*DIP, 13*DIP), size(13*DIP, 11*DIP)), Homogeneous, Stretch, (26*2 + 11) * DIP, 13 * DIP},
	}

	for i, v := range cases {
		in := hboxElement{
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
