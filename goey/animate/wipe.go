package animate

import (
	"guitest/goey/base"
)

var (
	wipeKind = base.NewKind("guitest/goey/animate.Wipe")
)

// Wipe provides an animation when changing its child widget.
type Wipe struct {
	Child base.Widget
	Level int
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Wipe) Kind() *base.Kind {
	return &wipeKind
}

// Mount creates a button control in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Wipe) Mount(parent base.Control) (base.Element, error) {
	child, err := base.Mount(parent, w.Child)
	if err != nil {
		return nil, err
	}

	return &wipeElement{
		child:  child,
		parent: parent,
		level:  w.Level,
	}, nil
}

type wipeElement struct {
	child    base.Element
	oldChild base.Element
	parent   base.Control
	level    int
	bounds   base.Rectangle
	ease     EaseLength
}

func (w *wipeElement) AnimateFrame(time Time) bool {
	dv := w.ease.Value(time)
	dx := w.bounds.Dx()

	bounds := w.bounds.Add(base.Point{dv, 0})
	w.child.SetBounds(bounds)
	if dv < 0 {
		bounds = w.bounds.Add(base.Point{dx + dv, 0})
	} else {
		bounds = w.bounds.Add(base.Point{-dx + dv, 0})
	}
	w.oldChild.SetBounds(bounds)

	if w.ease.Done(time) {
		w.oldChild.Close()
		w.oldChild = nil
		w.paint()
		return false
	}

	w.paint()

	return true
}

func (w *wipeElement) Close() {
	w.child.Close()
	if w.oldChild != nil {
		w.oldChild.Close()
		w.oldChild = nil
	}
}

func (*wipeElement) Kind() *base.Kind {
	return &wipeKind
}

func (w *wipeElement) Layout(bc base.Constraints) base.Size {
	return w.child.Layout(bc)
}

func (w *wipeElement) MinIntrinsicHeight(width base.Length) base.Length {
	return w.child.MinIntrinsicHeight(width)
}

func (w *wipeElement) MinIntrinsicWidth(height base.Length) base.Length {
	return w.child.MinIntrinsicWidth(height)
}

func (w *wipeElement) SetBounds(bounds base.Rectangle) {
	if bounds == w.bounds {
		return
	}

	w.bounds = bounds
	w.child.SetBounds(bounds)
}

func (w *wipeElement) updateProps(data *Wipe) error {
	if w.level == data.Level {
		child, err := base.DiffChild(w.parent, w.child, data.Child)
		if err != nil {
			return err
		}
		w.child = child
		return nil
	}

	child, err := base.Mount(w.parent, data.Child)
	if err != nil {
		return err
	}

	if w.oldChild != nil {
		w.oldChild.Close()
		w.oldChild = nil
	}

	w.oldChild = w.child
	w.child = child
	if w.level > data.Level {
		w.ease = NewEaseLength(0, w.bounds.Dx())
	} else {
		w.ease = NewEaseLength(0, -w.bounds.Dx())
	}
	child.Layout(base.Tight(base.Size{w.bounds.Dx(), w.bounds.Dy()}))
	child.SetBounds(w.bounds.Add(base.Point{w.ease.ca, 0}))

	w.level = data.Level
	AddAnimation(w)
	return nil
}

func (w *wipeElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*Wipe))
}
