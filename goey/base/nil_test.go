package base

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func ExampleMount() {
	// This won't work in real code, as the zero value for a control is not
	// generally useable.
	parent := Control{}

	// It is okay to mount a nil widget.
	elem, err := Mount(parent, nil)
	if err != nil {
		panic("Unexpected error!")
	}
	defer elem.Close()
	fmt.Println("The value of elem is nil...", elem == nil)
	fmt.Println("The kind of elem is...", elem.Kind())

	// Output:
	// The value of elem is nil... false
	// The kind of elem is... guitest/goey/base.nil
}

func TestMount(t *testing.T) {
	kind1 := NewKind("guitest/goey/base.Mock1")
	kind2 := NewKind("guitest/goey/base.Mock2")
	err1 := errors.New("fake error 1 for mounting widget")
	err2 := errors.New("fake error 2 for mounting widget")

	cases := []struct {
		in  Widget
		out Element
		err error
	}{
		{nil, (*nilElement)(nil), nil},
		{&mock{kind: &kind1}, &mockElement{kind: &kind1}, nil},
		{&mock{kind: &kind1, Prop: 3}, &mockElement{kind: &kind1, Prop: 3}, nil},
		{&mock{kind: &kind2}, &mockElement{kind: &kind2}, nil},
		{&mock{kind: &kind2, Prop: 13}, &mockElement{kind: &kind2, Prop: 13}, nil},
		{&mock{kind: &kind1, err: err1}, nil, err1},
		{&mock{kind: &kind1, err: err2}, nil, err2},
	}

	for i, v := range cases {
		out, err := Mount(Control{}, v.in)
		if err != v.err {
			t.Errorf("Case %d: Returned error does not match, got %v, want %v", i, err, v.err)
		}
		if !reflect.DeepEqual(out, v.out) {
			t.Errorf("Case %d: Returned element does not match, got %v, want %v", i, out, v.out)
		}
	}
}
