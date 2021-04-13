package main

import (
	"guitest/goey/base"
)

var (
	minsizedboxKind = base.NewKind("guitest/goey/example/feettometer.MinSizedBox")
)

// MinSizedBox is a custom layout widget that sizes its child widget according
// to the MinIntrinsicWidth and MinIntrinsicHeight of that widget.
type MinSizedBox struct {
	Child base.Widget // Child widget.
}

// Kind returns the concrete type for use in the.Widget interface.
// Users should not need to use this method directly.
func (*MinSizedBox) Kind() *base.Kind {
	return &minsizedboxKind
}

// Mount creates a this child in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *MinSizedBox) Mount(parent base.Control) (base.Element, error) {
	// Mount the child
	child, err := w.Child.Mount(parent)
	if err != nil {
		return nil, err
	}

	return &minsizedboxElement{
		parent: parent,
		child:  child,
	}, nil
}

type minsizedboxElement struct {
	parent    base.Control
	child     base.Element
	childSize base.Size
}

func (w *minsizedboxElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
}

func (*minsizedboxElement) Kind() *base.Kind {
	return &minsizedboxKind
}

func (w *minsizedboxElement) Layout(bc base.Constraints) base.Size {
	if w.child == nil {
		return bc.Constrain(base.Size{})
	}

	width := w.child.MinIntrinsicWidth(0)
	height := w.child.MinIntrinsicHeight(width)

	size := bc.Constrain(base.Size{width, height})
	return w.child.Layout(base.Tight(size))
}

func (w *minsizedboxElement) MinIntrinsicHeight(width base.Length) base.Length {
	if w.child == nil {
		return 0
	}

	return w.child.MinIntrinsicHeight(width)
}

func (w *minsizedboxElement) MinIntrinsicWidth(height base.Length) base.Length {
	if w.child == nil {
		return 0
	}

	return w.child.MinIntrinsicWidth(height)
}

func (w *minsizedboxElement) SetBounds(bounds base.Rectangle) {
	if w.child != nil {
		w.child.SetBounds(bounds)
	}
}

func (w *minsizedboxElement) updateProps(data *MinSizedBox) (err error) {
	w.child, err = base.DiffChild(w.parent, w.child, data.Child)
	w.childSize = base.Size{}
	return err
}

func (w *minsizedboxElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*MinSizedBox))
}
