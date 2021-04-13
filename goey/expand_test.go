package goey

import (
	"guitest/goey/base"
	"guitest/goey/mock"
	"errors"
	"testing"
)

type Boundser interface {
	Bounds() base.Rectangle
}

func (w *expandElement) Bounds() base.Rectangle {
	return w.child.(Boundser).Bounds()
}

func (w *expandElement) Props() base.Widget {
	child := base.Widget(nil)
	if w.child != nil {
		child = w.child.(Proper).Props()
	}

	return &Expand{
		Factor: w.factor,
		Child:  child,
	}
}

func TestExpandMount(t *testing.T) {
	testingMountWidgets(t,
		&Expand{},
		&Expand{Child: &mock.Widget{}},
	)

	// This should mount with an error.
	err := errors.New("Mock error 1")
	testingMountWidgetsFail(t, err,
		&Expand{Child: &mock.Widget{Err: err}},
	)
}

func TestExpandClose(t *testing.T) {
	testingCloseWidgets(t,
		&Expand{},
		&Expand{Child: &mock.Widget{}},
	)
}

func TestExpandUpdateProps(t *testing.T) {
	child := mock.Widget{}

	testingUpdateWidgets(t, []base.Widget{
		&Expand{},
		&Expand{Child: &child},
		&Expand{Child: &child, Factor: 1},
	}, []base.Widget{
		&Expand{Child: &child, Factor: 2},
		&Expand{},
		&Expand{Child: &child, Factor: 1},
	})
}

func TestExpandLayout(t *testing.T) {
	children := []base.Element{
		mock.New(base.Size{10 * DIP, 10 * DIP}),
		&expandElement{child: mock.New(base.Size{10 * DIP, 10 * DIP})},
		mock.New(base.Size{20 * DIP, 20 * DIP}),
	}

	cases := []struct {
		children    []base.Element
		alignMain   MainAxisAlign
		alignCross  CrossAxisAlign
		constraints base.Constraints
		size        base.Size
		bounds      []base.Rectangle
	}{
		{children, MainStart, Stretch, base.TightWidth(40 * DIP), base.Size{40 * DIP, 62 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 40*DIP, 10*DIP), base.Rect(0, 21*DIP, 40*DIP, 31*DIP), base.Rect(0, 42*DIP, 40*DIP, 62*DIP),
		}},
		{children, MainCenter, Stretch, base.TightWidth(40 * DIP), base.Size{40 * DIP, 62 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 40*DIP, 10*DIP), base.Rect(0, 21*DIP, 40*DIP, 31*DIP), base.Rect(0, 42*DIP, 40*DIP, 62*DIP),
		}},
		{children, MainEnd, Stretch, base.TightWidth(40 * DIP), base.Size{40 * DIP, 62 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 40*DIP, 10*DIP), base.Rect(0, 21*DIP, 40*DIP, 31*DIP), base.Rect(0, 42*DIP, 40*DIP, 62*DIP),
		}},
		{children, MainStart, CrossStart, base.TightWidth(40 * DIP), base.Size{40 * DIP, 62 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 10*DIP, 10*DIP), base.Rect(0, 21*DIP, 10*DIP, 31*DIP), base.Rect(0, 42*DIP, 20*DIP, 62*DIP),
		}},
		{children, MainStart, CrossCenter, base.TightWidth(40 * DIP), base.Size{40 * DIP, 62 * DIP}, []base.Rectangle{
			base.Rect(15*DIP, 0, 25*DIP, 10*DIP), base.Rect(15*DIP, 21*DIP, 25*DIP, 31*DIP), base.Rect(10*DIP, 42*DIP, 30*DIP, 62*DIP),
		}},
		{children, MainStart, CrossEnd, base.TightWidth(40 * DIP), base.Size{40 * DIP, 62 * DIP}, []base.Rectangle{
			base.Rect(30*DIP, 0, 40*DIP, 10*DIP), base.Rect(30*DIP, 21*DIP, 40*DIP, 31*DIP), base.Rect(20*DIP, 42*DIP, 40*DIP, 62*DIP),
		}},
		{children, MainStart, Stretch, base.Tight(base.Size{40 * DIP, 80 * DIP}), base.Size{40 * DIP, 80 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 40*DIP, 10*DIP), base.Rect(0, 21*DIP, 40*DIP, (31+18)*DIP), base.Rect(0, 60*DIP, 40*DIP, 80*DIP),
		}},
		{children, MainStart, CrossStart, base.Tight(base.Size{40 * DIP, 80 * DIP}), base.Size{40 * DIP, 80 * DIP}, []base.Rectangle{
			base.Rect(0, 0, 10*DIP, 10*DIP), base.Rect(0, 21*DIP, 10*DIP, (31+18)*DIP), base.Rect(0, 60*DIP, 20*DIP, 80*DIP),
		}},
		{children, MainStart, CrossCenter, base.Tight(base.Size{40 * DIP, 80 * DIP}), base.Size{40 * DIP, 80 * DIP}, []base.Rectangle{
			base.Rect(15*DIP, 0, 25*DIP, 10*DIP), base.Rect(15*DIP, 21*DIP, 25*DIP, (31+18)*DIP), base.Rect(10*DIP, 60*DIP, 30*DIP, 80*DIP),
		}},
		{children, MainStart, CrossEnd, base.Tight(base.Size{40 * DIP, 80 * DIP}), base.Size{40 * DIP, 80 * DIP}, []base.Rectangle{
			base.Rect(30*DIP, 0, 40*DIP, 10*DIP), base.Rect(30*DIP, 21*DIP, 40*DIP, (31+18)*DIP), base.Rect(20*DIP, 60*DIP, 40*DIP, 80*DIP),
		}},
	}

	for i, v := range cases {
		in := vboxElement{
			children:     v.children,
			alignMain:    v.alignMain,
			alignCross:   v.alignCross,
			childrenInfo: make([]boxElementInfo, len(v.children)),
		}
		in.childrenInfo, in.totalFlex = updateFlex(in.children, in.alignMain, in.childrenInfo)

		size := in.Layout(v.constraints)
		if size != v.size {
			t.Errorf("Incorrect size on case %d, got %s, want %s", i, size, v.size)
		}
		in.SetBounds(base.Rect(0, 0, size.Width, size.Height))
		for j, u := range v.bounds {
			if got := v.children[j].(Boundser).Bounds(); got != u {
				t.Errorf("Incorrect bounds case %d-%d, got %s, want %s", i, j, got, u)
			}
		}
	}
}
