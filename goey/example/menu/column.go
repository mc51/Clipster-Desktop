package main

import (
	"guitest/goey/base"
)

var (
	columnKind = base.NewKind("guitest/goey/example/menu.Column")
)

// Column describes a layout widget that arranges its child widgets into a several columns.
// If there is sufficient width, the columns will be arranged side-by-side.  Otherwise, the
// columns will be arranged vertically.
type Column struct {
	Children []base.Widget
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Column) Kind() *base.Kind {
	return &columnKind
}

// Mount creates a horizontal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Column) Mount(parent base.Control) (base.Element, error) {
	c := make([]base.Element, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			return nil, err
		}
		c = append(c, mountedChild)
	}

	return &columnElement{
		parent:     parent,
		children:   c,
		rowHeights: []base.Length{},
	}, nil
}

type columnElement struct {
	parent   base.Control
	children []base.Element

	minWidth   base.Length
	rowHeights []base.Length
}

func (w *columnElement) Close() {
	base.CloseElements(w.children)
	w.children = nil
}

func (*columnElement) Kind() *base.Kind {
	return &columnKind
}

func (w *columnElement) Layout(bc base.Constraints) base.Size {
	const gap = 11 * base.DIP

	if len(w.children) == 0 {
		return bc.Constrain(base.Size{})
	}

	if w.minWidth == 0 {
		w.MinIntrinsicWidth(base.Inf)
	}

	width := bc.Max.Width
	if !bc.HasBoundedWidth() {
		width = bc.ConstrainWidth(w.minWidth)
	}

	columns := int((width + gap) / (w.minWidth + gap))
	cbc := base.TightWidth((width + gap).Scale(1, columns))

	height := -gap
	rowHeight := base.Length(0)
	w.rowHeights = w.rowHeights[:0]
	for i, v := range w.children {
		tmp := v.Layout(cbc)
		if tmp.Height > rowHeight {
			rowHeight = tmp.Height
		}
		if (i+1)%columns == 0 {
			height += rowHeight + gap
			w.rowHeights = append(w.rowHeights, rowHeight)
			rowHeight = 0
		}
	}
	if rowHeight > 0 {
		height += rowHeight + gap
		w.rowHeights = append(w.rowHeights, rowHeight)
	}

	return base.Size{width, bc.ConstrainHeight(height)}
}

func (w *columnElement) MinIntrinsicWidth(height base.Length) base.Length {
	if len(w.children) == 0 {
		return 0
	}

	if height != base.Inf {
		panic("not implemented")
	}

	minWidth := w.children[0].MinIntrinsicWidth(base.Inf)
	for _, v := range w.children[1:] {
		tmp := v.MinIntrinsicWidth(base.Inf)
		if tmp > minWidth {
			minWidth = tmp
		}
	}
	w.minWidth = minWidth
	return minWidth
}

func (w *columnElement) MinIntrinsicHeight(width base.Length) base.Length {
	if len(w.children) == 0 {
		return 0
	}

	if width != base.Inf {
		panic("not implemented")
	}

	minHeight := w.children[0].MinIntrinsicHeight(base.Inf)
	for _, v := range w.children[1:] {
		tmp := v.MinIntrinsicHeight(base.Inf)
		if tmp > minHeight {
			minHeight = tmp
		}
	}
	return minHeight
}

func (w *columnElement) SetBounds(bounds base.Rectangle) {
	const gap = 11 * base.DIP

	if len(w.children) == 0 {
		return
	}

	width := bounds.Dx()
	columns := int((width + gap) / (w.minWidth + gap))
	itemWidth := (width+gap).Scale(1, columns) - gap

	posX := bounds.Min.X
	posY := bounds.Min.Y
	rowHeights := w.rowHeights
	for i, v := range w.children {
		v.SetBounds(base.Rect(posX, posY, posX+itemWidth, posY+rowHeights[0]))
		posX += itemWidth + gap
		if (i+1)%columns == 0 {
			posY += rowHeights[0] + gap
			posX = bounds.Min.X
			rowHeights = rowHeights[1:]
		}
	}
}

func (w *columnElement) updateProps(data *Column) (err error) {
	w.minWidth = 0
	w.children, err = base.DiffChildren(w.parent, w.children, data.Children)
	return err
}

func (w *columnElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*Column))
}
