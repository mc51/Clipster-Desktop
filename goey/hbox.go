package goey

import (
	"clipster/goey/base"
)

var (
	hboxKind = base.NewKind("clipster/goey.HBox")
)

// HBox describes a layout widget that arranges its child widgets into a row.
// Children are positioned in order from the left towards the right.  The main
// axis for alignment is therefore horizontal, with the cross axis for alignment is vertical.
//
// The size of the box will try to set a width sufficient to contain all of its
// children.  Extra space will be distributed according to the value of
// AlignMain.  Subject to the box constraints during layout, the height should
// match the largest minimum height of the child widgets.
type HBox struct {
	AlignMain  MainAxisAlign  // Control distribution of excess horizontal space when positioning children.
	AlignCross CrossAxisAlign // Control distribution of excess vertical space when positioning children.
	Children   []base.Widget  // Children.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*HBox) Kind() *base.Kind {
	return &hboxKind
}

// Mount creates a horizontal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *HBox) Mount(parent base.Control) (base.Element, error) {
	c := make([]base.Element, 0, len(w.Children))

	// Mount all of the children
	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			base.CloseElements(c)
			return nil, err
		}
		c = append(c, mountedChild)
	}

	// Record the flex factor for all children
	ci, totalFlex := updateFlex(c, w.AlignMain, nil)

	return &hboxElement{
		parent:       parent,
		children:     c,
		alignMain:    w.AlignMain,
		alignCross:   w.AlignCross,
		childrenInfo: ci,
		totalFlex:    totalFlex,
	}, nil
}

type hboxElement struct {
	parent     base.Control
	children   []base.Element
	alignMain  MainAxisAlign
	alignCross CrossAxisAlign

	childrenInfo []boxElementInfo
	totalWidth   base.Length
	totalFlex    int
}

type boxElementInfo struct {
	size base.Size
	flex int
}

func (w *hboxElement) Close() {
	base.CloseElements(w.children)
	w.children = nil
	w.childrenInfo = nil
}

func (*hboxElement) Kind() *base.Kind {
	return &hboxKind
}

func (w *hboxElement) Layout(bc base.Constraints) base.Size {
	if len(w.children) == 0 {
		w.totalWidth = 0
		return bc.Constrain(base.Size{})
	}

	// Determine the constraints for layout of child elements.
	cbc := bc
	if w.alignMain == Homogeneous {
		count := len(w.children)
		gap := calculateHGap(nil, nil)
		cbc.TightenWidth(cbc.Max.Width.Scale(1, count) - gap.Scale(count-1, count))
	} else {
		cbc.Min.Width = 0
		cbc.Max.Width = base.Inf
	}
	if w.alignCross == Stretch {
		if cbc.HasBoundedHeight() {
			cbc = cbc.TightenHeight(cbc.Max.Height)
		} else {
			cbc = cbc.TightenHeight(w.MinIntrinsicHeight(base.Inf))
		}
	} else {
		cbc = cbc.LoosenHeight()
	}

	width := base.Length(0)
	height := base.Length(0)
	previous := base.Element(nil)
	for i, v := range w.children {
		// Determine what gap needs to be inserted between the elements.
		if i > 0 {
			if w.alignMain.IsPacked() {
				width += calculateHGap(previous, v)
			} else {
				width += calculateHGap(nil, nil)
			}
		}
		previous = v

		// Perform layout of the element.  Track impact on width and height.
		size := v.Layout(cbc)
		w.childrenInfo[i].size = size
		width += size.Width
		height = max(height, size.Height)
	}
	w.totalWidth = width

	// Need to adjust height to any widgets that have flex
	if w.totalFlex > 0 {
		extraWidth := base.Length(0)
		if bc.HasBoundedWidth() && bc.Max.Width > w.totalWidth {
			extraWidth = bc.Max.Width - w.totalWidth
		} else if bc.Min.Width > w.totalWidth {
			extraWidth = bc.Min.Width - w.totalWidth
		}

		if extraWidth > 0 {
			for i, v := range w.childrenInfo {
				if v.flex > 0 {
					oldWidth := v.size.Width
					fbc := cbc.TightenWidth(v.size.Width + extraWidth.Scale(v.flex, w.totalFlex))
					size := w.children[i].Layout(fbc)
					w.childrenInfo[i].size = size
					w.totalWidth += size.Width - oldWidth
				}
			}
		}
	}

	if w.alignCross == Stretch {
		return bc.Constrain(base.Size{width, cbc.Min.Height})
	}
	return bc.Constrain(base.Size{width, height})
}

func (w *hboxElement) MinIntrinsicHeight(width base.Length) base.Length {
	if len(w.children) == 0 {
		return 0
	}

	if w.alignMain == Homogeneous {
		width = guardInf(width, width.Scale(1, len(w.children)))
		size := w.children[0].MinIntrinsicHeight(width)
		for _, v := range w.children[1:] {
			size = max(size, v.MinIntrinsicHeight(width))
		}
		return size
	}

	size := w.children[0].MinIntrinsicHeight(base.Inf)
	for _, v := range w.children[1:] {
		size = max(size, v.MinIntrinsicHeight(base.Inf))
	}
	return size
}

func (w *hboxElement) MinIntrinsicWidth(height base.Length) base.Length {
	if len(w.children) == 0 {
		return 0
	}

	size := w.children[0].MinIntrinsicWidth(height)
	if w.alignMain.IsPacked() {
		previous := w.children[0]
		for _, v := range w.children[1:] {
			// Add the preferred gap between this pair of widgets
			size += calculateHGap(previous, v)
			previous = v
			// Find minimum size for this widget, and update
			size += v.MinIntrinsicWidth(height)
		}
		return size
	}

	if w.alignMain == Homogeneous {
		for _, v := range w.children[1:] {
			size = max(size, v.MinIntrinsicWidth(height))
		}

		// Add a minimum gap between the controls.
		size = size.Scale(len(w.children), 1) + calculateHGap(nil, nil).Scale(len(w.children)-1, 1)
		return size
	}

	for _, v := range w.children[1:] {
		size += v.MinIntrinsicWidth(height)
	}

	// Add a minimum gap between the controls.
	if w.alignMain == SpaceBetween {
		size += calculateHGap(nil, nil).Scale(len(w.children)-1, 1)
	} else {
		size += calculateHGap(nil, nil).Scale(len(w.children)+1, 1)
	}

	return size
}

func (w *hboxElement) SetBounds(bounds base.Rectangle) {
	if len(w.children) == 0 {
		return
	}

	if w.alignMain == Homogeneous {
		gap := calculateHGap(nil, nil)
		dx := bounds.Dx() + gap
		count := len(w.children)

		for i, v := range w.children {
			x1 := bounds.Min.X + dx.Scale(i, count)
			x2 := bounds.Min.X + dx.Scale(i+1, count) - gap
			w.setBoundsForChild(i, v, x1, bounds.Min.Y, x2, bounds.Max.Y)
		}
		return
	}

	// Adjust the bounds so that the minimum Y handles vertical alignment
	// of the controls.  We also calculate 'extraGap' which will adjust
	// spacing of the controls for non-packed alignments.
	extraGap := base.Length(0)
	if w.totalFlex == 0 {
		switch w.alignMain {
		case MainStart:
			// Do nothing
		case MainCenter:
			bounds.Min.X += (bounds.Dx() - w.totalWidth) / 2
		case MainEnd:
			bounds.Min.X = bounds.Max.X - w.totalWidth
		case SpaceAround:
			extraGap = (bounds.Dx() - w.totalWidth).Scale(1, len(w.children)+1)
			bounds.Min.X += extraGap
			extraGap += calculateHGap(nil, nil)
		case SpaceBetween:
			if len(w.children) > 1 {
				extraGap = (bounds.Dx() - w.totalWidth).Scale(1, len(w.children)-1)
				extraGap += calculateHGap(nil, nil)
			} else {
				// There are no controls between which to put the extra space.
				// The following essentially convert SpaceBetween to SpaceAround
				bounds.Min.X += (bounds.Dx() - w.totalWidth) / 2
			}
		}
	}

	// Position all of the child controls.
	posX := bounds.Min.X
	previous := base.Element(nil)
	for i, v := range w.children {
		if w.alignMain.IsPacked() {
			if i > 0 {
				posX += calculateHGap(previous, v)
			}
			previous = v
		}

		dx := w.childrenInfo[i].size.Width
		w.setBoundsForChild(i, v, posX, bounds.Min.Y, posX+dx, bounds.Max.Y)
		posX += dx + extraGap
	}
}

func (w *hboxElement) setBoundsForChild(i int, v base.Element, posX, posY, posX2, posY2 base.Length) {
	dy := w.childrenInfo[i].size.Height
	switch w.alignCross {
	case CrossStart:
		v.SetBounds(base.Rectangle{
			base.Point{posX, posY},
			base.Point{posX2, posY + dy},
		})
	case CrossCenter:
		v.SetBounds(base.Rectangle{
			base.Point{posX, posY + (posY2-posY-dy)/2},
			base.Point{posX2, posY + (posY2-posY+dy)/2},
		})
	case CrossEnd:
		v.SetBounds(base.Rectangle{
			base.Point{posX, posY2 - dy},
			base.Point{posX2, posY2},
		})
	case Stretch:
		v.SetBounds(base.Rectangle{
			base.Point{posX, posY},
			base.Point{posX2, posY2},
		})
	}
}

func updateFlex(c []base.Element, alignMain MainAxisAlign, clientInfo []boxElementInfo) ([]boxElementInfo, int) {
	if len(c) <= cap(clientInfo) {
		clientInfo = clientInfo[:len(c)]
	} else {
		clientInfo = make([]boxElementInfo, len(c))
	}

	totalFlex := 0
	for i, v := range c {
		if elem, ok := v.(*expandElement); ok {
			clientInfo[i].flex = elem.factor + 1
			totalFlex += elem.factor + 1
		}
	}
	if alignMain == Homogeneous {
		totalFlex = 0
	}
	return clientInfo, totalFlex
}

func (w *hboxElement) updateProps(data *HBox) (err error) {
	// Update properties
	w.alignMain = data.AlignMain
	w.alignCross = data.AlignCross
	w.children, err = base.DiffChildren(w.parent, w.children, data.Children)
	// Clear cached values
	w.childrenInfo, w.totalFlex = updateFlex(w.children, w.alignMain, w.childrenInfo)
	w.totalWidth = 0
	return err
}

func (w *hboxElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*HBox))
}
