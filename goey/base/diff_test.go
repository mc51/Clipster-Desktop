package base

import (
	"errors"
	"reflect"
	"testing"
)

type mock struct {
	kind *Kind
	err  error
	Prop int
}

func (m *mock) Kind() *Kind {
	return m.kind
}

func (m *mock) Mount(parent Control) (Element, error) {
	// Check if the mock widget is supposed to fail with an error when mounted.
	if m.err != nil {
		return nil, m.err
	}

	// Create the mock element.
	return &mockElement{
		kind: m.kind,
		Prop: m.Prop,
	}, nil
}

type mockElement struct {
	kind   *Kind
	err    error
	Closed bool
	Prop   int
}

func (m *mockElement) Close() {
	m.Closed = true
}

func (m *mockElement) Kind() *Kind {
	return m.kind
}

func (m *mockElement) Layout(bc Constraints) Size {
	return bc.Constrain(Size{})
}
func (m *mockElement) MinIntrinsicHeight(width Length) Length {
	return 0
}

func (m *mockElement) MinIntrinsicWidth(height Length) Length {
	return 0
}

func (m *mockElement) SetBounds(bounds Rectangle) {

}

func (m *mockElement) updateProps(data *mock) error {
	if m.kind != data.kind {
		panic("Mismatched kinds")
	}
	if m.err != nil {
		return m.err
	}
	m.Prop = data.Prop
	return nil
}

func (m *mockElement) UpdateProps(data Widget) error {
	return m.updateProps(data.(*mock))
}

func TestCloseElements(t *testing.T) {
	kind := NewKind("guitest/goey/base.Mock")

	// Check for no panic on nil or zero-length list
	CloseElements(nil)
	CloseElements([]Element{})

	for _, v := range []int{1, 2, 3, 4, 8, 16} {
		elem := make([]Element, 0, v)
		for i := 0; i < v; i++ {
			elem = append(elem, &mockElement{kind: &kind})
		}

		CloseElements(elem)

		for _, v := range elem {
			if !v.(*mockElement).Closed {
				t.Errorf("Failed to close element")
			}
		}
	}
}

func TestDiffChild(t *testing.T) {
	kind1 := NewKind("guitest/goey/base.Mock1")
	kind2 := NewKind("guitest/goey/base.Mock2")
	err1 := errors.New("fake error 1 for mounting widget")
	err2 := errors.New("fake error 2 for mounting widget")

	cases := []struct {
		lhs       Element
		rhs       Widget
		out       Element
		err       error
		lhsClosed bool
	}{
		// Trivial case
		{nil, nil, (*nilElement)(nil), nil, false},
		{(*nilElement)(nil), nil, (*nilElement)(nil), nil, false},
		// Mount
		{nil, &mock{kind: &kind1}, &mockElement{kind: &kind1}, nil, false},
		{nil, &mock{kind: &kind1, Prop: 3}, &mockElement{kind: &kind1, Prop: 3}, nil, false},
		{nil, &mock{kind: &kind2}, &mockElement{kind: &kind2}, nil, false},
		{nil, &mock{kind: &kind2, Prop: 13}, &mockElement{kind: &kind2, Prop: 13}, nil, false},
		// Remove existing element
		{&mockElement{kind: &kind1}, nil, (*nilElement)(nil), nil, true},
		// Replace existing element
		{&mockElement{kind: &kind1, Prop: 3}, &mock{kind: &kind2, Prop: 13}, &mockElement{kind: &kind2, Prop: 13}, nil, true},
		// Update existing element
		{&mockElement{kind: &kind1, Prop: 3}, &mock{kind: &kind1, Prop: 13}, &mockElement{kind: &kind1, Prop: 13}, nil, false},
		// Fail to mount
		{nil, &mock{kind: &kind1, err: err1}, nil, err1, false},
		{nil, &mock{kind: &kind1, err: err2}, nil, err2, false},
		// Fail to replace
		{&mockElement{kind: &kind1}, &mock{kind: &kind2, err: err1}, &mockElement{kind: &kind1}, err1, false},
		{&mockElement{kind: &kind1}, &mock{kind: &kind2, err: err2}, &mockElement{kind: &kind1}, err2, false},
		// Fail to update
		{&mockElement{kind: &kind1, err: err1}, &mock{kind: &kind1, Prop: 4}, &mockElement{kind: &kind1, err: err1}, err1, false},
	}

	for i, v := range cases {
		out, err := DiffChild(Control{}, v.lhs, v.rhs)
		if err != v.err {
			if v.err == nil {
				t.Errorf("Case %d: Unexpected error during DiffChild, %s", i, err)
			} else {
				t.Errorf("Case %d: Returned error does not match, got %v, want %v", i, err, v.err)
			}
		}
		if !reflect.DeepEqual(out, v.out) {
			t.Errorf("Case %d: Returned element does not match, got %v, want %v", i, out, v.out)
		}
		if v.lhsClosed && !v.lhs.(*mockElement).Closed {
			t.Errorf("Case %d: Failed to close lhs", i)
		}
	}
}

func TestDiffChildren(t *testing.T) {
	kind1 := NewKind("guitest/goey/base.Mock1")
	kind2 := NewKind("guitest/goey/base.Mock2")
	err1 := errors.New("fake error 1 for mounting widget")
	err2 := errors.New("fake error 2 for mounting widget")

	cases := []struct {
		lhs       []Element
		rhs       []Widget
		out       []Element
		err       error
		lhsClosed bool
	}{
		// Trivial case
		{nil, nil, nil, nil, false},
		// Mount
		{nil, []Widget{&mock{kind: &kind1}}, []Element{&mockElement{kind: &kind1}}, nil, false},
		{nil, []Widget{&mock{kind: &kind1, Prop: 3}}, []Element{&mockElement{kind: &kind1, Prop: 3}}, nil, false},
		{nil, []Widget{&mock{kind: &kind2}}, []Element{&mockElement{kind: &kind2}}, nil, false},
		{nil, []Widget{&mock{kind: &kind2, Prop: 13}}, []Element{&mockElement{kind: &kind2, Prop: 13}}, nil, false},
		{nil, []Widget{&mock{kind: &kind1, err: err1}}, nil, err1, false},
		{nil, []Widget{&mock{kind: &kind1, err: err2}}, nil, err2, false},
		// Remove existing elements
		{[]Element{}, nil, nil, nil, true},
		{[]Element{&mockElement{kind: &kind1}}, nil, nil, nil, true},
		{[]Element{&mockElement{kind: &kind2}}, nil, nil, nil, true},
		// Update existing element
		{
			[]Element{&mockElement{kind: &kind1}},
			[]Widget{&mock{kind: &kind1, Prop: 15}},
			[]Element{&mockElement{kind: &kind1, Prop: 15}},
			nil, false,
		},
		{
			[]Element{&mockElement{kind: &kind2}},
			[]Widget{&mock{kind: &kind2, Prop: 16}},
			[]Element{&mockElement{kind: &kind2, Prop: 16}},
			nil, false,
		},
		// Remove extra elements
		{
			[]Element{&mockElement{kind: &kind1}, &mockElement{kind: &kind2, Prop: 17}},
			[]Widget{&mock{kind: &kind1, Prop: 15}},
			[]Element{&mockElement{kind: &kind1, Prop: 15}},
			nil, false,
		},
		// Mount new elements
		{
			[]Element{&mockElement{kind: &kind1}},
			[]Widget{&mock{kind: &kind1, Prop: 15}, &mock{kind: &kind2, Prop: 17}},
			[]Element{&mockElement{kind: &kind1, Prop: 15}, &mockElement{kind: &kind2, Prop: 17}},
			nil, false,
		},
		// Replace existing element
		{
			[]Element{&mockElement{kind: &kind1, Prop: 123}},
			[]Widget{&mock{kind: &kind2}},
			[]Element{&mockElement{kind: &kind2}},
			nil, true,
		},
		// Fail to replace existing element
		{
			[]Element{&mockElement{kind: &kind1}},
			[]Widget{&mock{kind: &kind2, err: err1}},
			[]Element{&mockElement{kind: &kind1}},
			err1, false,
		},
		// Fail to add new element
		{
			[]Element{&mockElement{kind: &kind1}},
			[]Widget{&mock{kind: &kind1}, &mock{kind: &kind1, err: err1}},
			[]Element{&mockElement{kind: &kind1}},
			err1, false,
		},
		// Fail to update existing element
		{
			[]Element{&mockElement{kind: &kind1}, &mockElement{kind: &kind2, err: err1}},
			[]Widget{&mock{kind: &kind1}, &mock{kind: &kind2, Prop: 1}},
			[]Element{&mockElement{kind: &kind1}, &mockElement{kind: &kind2, err: err1}},
			err1, false,
		},
	}

	for i, v := range cases {
		out, err := DiffChildren(Control{}, append([]Element(nil), v.lhs...), v.rhs)
		if err != v.err {
			if v.err == nil {
				t.Errorf("Case %d: Unexpected error during DiffChildren, %s", i, err)
			} else {
				t.Errorf("Case %d: Returned error does not match, got %v, want %v", i, err, v.err)
			}
		}
		if !reflect.DeepEqual(out, v.out) {
			t.Errorf("Case %d: Returned element does not match, got %v, want %v", i, out, v.out)
		}
		if len(out) < len(v.lhs) {
			for j, v := range v.lhs[len(out):] {
				if !v.(*mockElement).Closed {
					t.Errorf("Case %d: Failed to close lhs[%d]", i, len(out)+j)
				}
			}
		}
		if v.lhsClosed {
			for j, v := range v.lhs {
				if !v.(*mockElement).Closed {
					t.Errorf("Case %d: Failed to close lhs[%d]", i, j)
				}
			}
		}
	}
}

func TestLayout(t *testing.T) {
	size1 := Size{96 * DIP, 2 * 96 * DIP}
	cases := []struct {
		in  Element
		bc  Constraints
		out Size
	}{
		{&nilElement{}, Tight(size1), size1},
		{&nilElement{}, Loose(size1), size1},
		{&nilElement{}, TightWidth(size1.Width), Size{size1.Width, 0}},
		{&nilElement{}, TightHeight(size1.Height), Size{0, size1.Height}},
		{&nilElement{}, Expand(), Size{}},
		{&mockElement{}, Tight(size1), size1},
		{&mockElement{}, Loose(size1), Size{}},
		{&mockElement{}, TightWidth(size1.Width), Size{size1.Width, 0}},
		{&mockElement{}, TightHeight(size1.Height), Size{0, size1.Height}},
		{&mockElement{}, Expand(), Size{}},
	}

	for i, v := range cases {
		out := v.in.Layout(v.bc)
		if out != v.out {
			t.Errorf("Case %d: Returned size does not match, got %v, want %v", i, out, v.out)
		}
	}
}
