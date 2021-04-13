// +build cocoa darwin,!gtk

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/cocoa"
)

type tabsElement struct {
	control  *cocoa.TabView
	value    int
	child    base.Element
	widgets  []TabItem
	insets   Insets
	onChange func(int)

	cachedBounds base.Rectangle
	cachedInsets base.Point
	cachedTabsW  base.Length
}

func (w *Tabs) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewTabView(parent.Handle)
	for _, v := range w.Children {
		control.AddItem(v.Caption)
	}

	child := base.Element(nil)
	if len(w.Children) > 0 {
		parent := base.Control{control.ContentView(w.Value)}
		child_, err := base.Mount(parent, w.Children[w.Value].Child)
		if err != nil {
			control.Close()
			return nil, err
		}
		child = child_
		control.SelectItem(w.Value)
	}

	retval := &tabsElement{
		control:  control,
		child:    child,
		value:    w.Value,
		widgets:  w.Children,
		insets:   w.Insets,
		onChange: w.OnChange,
	}
	control.SetOnChange(retval.OnChange)
	return retval, nil
}

func (w *tabsElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *tabsElement) contentInsets() base.Point {
	if w.cachedInsets.Y == 0 {
		x, y := w.control.ContentInsets()
		w.cachedInsets.X = base.FromPixelsX(x)
		w.cachedInsets.Y = base.FromPixelsY(y)
	}

	return w.cachedInsets
}

func (w *tabsElement) controlTabsMinWidth() base.Length {
	if w.cachedTabsW == 0 {
		w.cachedTabsW = 100 * base.DIP
	}
	return w.cachedTabsW
}

func (w *tabsElement) mountPage(page int) {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}

	parent := base.Control{w.control.ContentView(page)}
	child, _ := base.Mount(parent, w.widgets[page].Child)
	child.Layout(base.Tight(base.Size{
		w.cachedBounds.Dx(),
		w.cachedBounds.Dy(),
	}))
	child.SetBounds(w.cachedBounds)
	w.child = child
}

func (w *tabsElement) OnChange(page int) {
	if page != w.value {
		if w.onChange != nil {
			w.onChange(page)
		}
		if page != w.value {
			w.mountPage(page)
			w.value = page
		}
	}
}

func (w *tabsElement) Props() base.Widget {
	count := w.control.NumberOfItems()
	children := make([]TabItem, count)
	for i := 0; i < count; i++ {
		label := w.control.ItemAtIndex(i)
		children[i].Caption = label
		children[i].Child = w.widgets[i].Child
	}

	return &Tabs{
		Value:    w.value,
		Children: children,
		Insets:   w.insets,
		OnChange: w.onChange,
	}
}

func (w *tabsElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())

	if w.child != nil {
		// Determine the bounds for the child widget
		dx := bounds.Dx() - w.cachedInsets.X
		dy := bounds.Dy() - w.cachedInsets.Y

		bounds.Min.X = w.insets.Left
		bounds.Min.Y = w.insets.Top
		bounds.Max.X = dx - w.insets.Right
		bounds.Max.Y = dy - w.insets.Bottom

		// Update bounds for the child
		w.cachedBounds = bounds
		w.child.SetBounds(bounds)
	}
}
func (w *tabsElement) updateProps(data *Tabs) error {
	if len(w.widgets) > len(data.Children) {
		// Modify captions for existing tabs
		for i, v := range data.Children {
			w.control.SetItemAtIndex(i, v.Caption)
		}
		// Remove excess tabs
		for i := len(w.widgets); i > len(data.Children); i-- {
			w.control.RemoveItemAtIndex(i - 1)
		}
	} else {
		// Modify captions for existing tabs
		for i, v := range data.Children[:len(w.widgets)] {
			w.control.SetItemAtIndex(i, v.Caption)
		}
		// Append new tabs
		for _, v := range data.Children[len(w.widgets):] {
			w.control.AddItem(v.Caption)
		}
	}
	w.widgets = data.Children

	// Update the selected widget
	if data.Value != w.value {
		w.mountPage(data.Value)

		w.control.SelectItem(data.Value)
		w.value = data.Value
	} else {
		parent := base.Control{w.control.ContentView(w.value)}
		child, err := base.DiffChild(parent, w.child, data.Children[data.Value].Child)
		w.child = child
		if err != nil {
			return err
		}
	}

	return nil
}
