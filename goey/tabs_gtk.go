// +build gtk linux darwin freebsd openbsd

package goey

import (
	"bytes"

	"clipster/goey/base"
	"clipster/goey/internal/gtk"
)

type tabsElement struct {
	Control
	value    int
	child    base.Element
	widgets  []TabItem
	insets   Insets
	onChange func(int)

	cachedInsets base.Point
	cachedBounds base.Rectangle
	cachedTabsW  base.Length
}

func (w *Tabs) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountTabs(parent.Handle, w.Value, w.serializeItems(),
		w.OnChange != nil)

	child := base.Element(nil)
	if len(w.Children) > 0 {
		parent := gtk.TabsGetTabParent(control, w.Value)
		child_, err := base.Mount(base.Control{parent}, w.Children[w.Value].Child)
		if err != nil {
			gtk.WidgetClose(control)
			return nil, err
		}
		child = child_
	}

	retval := &tabsElement{
		Control:  Control{control},
		child:    child,
		value:    w.Value,
		widgets:  w.Children,
		insets:   w.Insets,
		onChange: w.OnChange,
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}

func (w *Tabs) serializeItems() string {
	buffer := bytes.Buffer{}

	for _, v := range w.Children {
		buffer.WriteString(v.Caption)
		buffer.WriteByte(0)
	}
	buffer.WriteByte(0)

	return buffer.String()
}

func (w *tabsElement) OnChange(page int) {
	if page != w.value {
		if w.onChange != nil {
			w.onChange(page)
		}
		if page != w.value {
			// Not clear how a handle at this point should be handled.
			// The widget is supposed to already be mounted, but we create and
			// remove controls when the tab is changed.
			// In practice, errors are very infrequent (never?).  GTK widgets
			// will never fail to mount.
			_ = w.mountPage(page)
			w.value = page
		}
	}
}

func (w *tabsElement) contentInsets() base.Point {
	if w.cachedInsets.Y == 0 {
		h1 := gtk.WidgetMinHeight(w.handle)
		// How should the offset between the notebook widget and the contained
		// page be measured?
		w.cachedInsets = base.Point{
			X: 0,
			Y: base.FromPixelsY(h1),
		}
	}

	return w.cachedInsets
}

func (w *tabsElement) controlTabsMinWidth() base.Length {
	if w.cachedTabsW == 0 {
		w1 := gtk.WidgetMinWidth(w.handle)
		w.cachedTabsW = base.FromPixelsX(w1)
	}
	return w.cachedTabsW
}

func (w *tabsElement) mountPage(page int) error {
	parent := gtk.TabsGetTabParent(w.handle, page)
	child, err := w.widgets[page].Child.Mount(base.Control{parent})
	if err != nil {
		return err
	}
	child.Layout(base.Tight(base.Size{
		Width:  w.cachedBounds.Dx(),
		Height: w.cachedBounds.Dy(),
	}))
	child.SetBounds(w.cachedBounds)

	if w.child != nil {
		w.child.Close()
	}
	w.child = child
	return nil
}

func (w *tabsElement) Props() base.Widget {
	count := gtk.TabsItemCount(w.handle)

	children := make([]TabItem, count)
	for i := 0; i < count; i++ {
		children[i].Caption = gtk.TabsItemCaption(w.handle, i)
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
	control := Control{w.handle}
	control.SetBounds(bounds)

	if w.child != nil {
		// Determine the bounds for the child widget
		insets := w.contentInsets()
		insets.X += w.insets.Left + w.insets.Right
		insets.Y += w.insets.Top + w.insets.Bottom
		bounds = base.Rectangle{
			Max: base.Point{bounds.Dx() - insets.X, bounds.Dy() - insets.Y},
		}

		// Offset
		offset := base.Point{w.insets.Left, w.insets.Top}
		bounds.Min = bounds.Min.Add(offset)
		bounds.Max = bounds.Max.Add(offset)

		// Update bounds for the child
		w.cachedBounds = bounds
		w.child.SetBounds(bounds)
	}
}

func (w *tabsElement) updateProps(data *Tabs) error {
	gtk.TabsUpdate(w.handle, data.Value, data.serializeItems(),
		data.OnChange != nil)
	w.widgets = data.Children

	// Update the selected widget
	if data.Value == w.value {
		err := w.mountPage(data.Value)
		if err != nil {
			return err
		}
	} else {
		parent := gtk.TabsGetTabParent(w.handle, data.Value)
		child, err := base.DiffChild(base.Control{parent}, w.child, data.Children[data.Value].Child)
		w.child = child
		if err != nil {
			return err
		}

		w.value = data.Value
	}

	return nil
}
