package goey

import (
	"guitest/goey/base"
)

var (
	tabsKind = base.NewKind("guitest/goey.Tabs")
)

// Tabs describes a widget that shows a tabs.
//
// The size of the control will match the size of the currently selected child
// element, although padding will added as required to provide space for the
// border and the tabs.  However, when the user switches tabs, a relayout of
// the entire window is not forced.
//
// When calling UpdateProps, setting Value to an integer less than zero will
// leave the currently selected tab unchanged.
type Tabs struct {
	Value    int       // Index of the selected tab
	Children []TabItem // Description of the tabs
	Insets   Insets    // Space between edge of element and the child element.

	OnChange func(int) // OnChange will be called whenever the user selects a different tab
}

// TabItem describes a tab for a Tab widget.
type TabItem struct {
	Caption string      // Text to describe the contents of this tab
	Child   base.Widget // Child widget for the tab
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Tabs) Kind() *base.Kind {
	return &tabsKind
}

// Mount creates a tabs control in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Tabs) Mount(parent base.Control) (base.Element, error) {
	// Ensure that the Value is a useable index.
	if w.Value < 0 {
		w.Value = 0
	}
	w.UpdateValue()
	// Forward to the platform-dependant code
	return w.mount(parent)
}

// UpdateValue ensures that the index for the currently selected tab is with
// the allowed range.
func (w *Tabs) UpdateValue() {
	if w.Value >= len(w.Children) {
		w.Value = len(w.Children) - 1
	}
}

func (*tabsElement) Kind() *base.Kind {
	return &tabsKind
}

func (w *tabsElement) Layout(bc base.Constraints) base.Size {
	insets := w.contentInsets()
	insets.X += w.insets.Left + w.insets.Right
	insets.Y += w.insets.Top + w.insets.Bottom

	if w.child == nil {
		return bc.Constrain(base.Size{
			Width:  insets.X,
			Height: insets.Y,
		})
	}

	size := w.child.Layout(bc.Inset(insets.X, insets.Y))
	return base.Size{
		Width:  size.Width + insets.X,
		Height: size.Height + insets.Y,
	}
}

func (w *tabsElement) MinIntrinsicHeight(width base.Length) base.Length {
	xInset := w.insets.Left + w.insets.Right
	yInset := w.insets.Top + w.insets.Bottom

	if w.child == nil {
		return yInset
	}

	if width == base.Inf {
		return w.child.MinIntrinsicHeight(base.Inf) + yInset
	}

	return w.child.MinIntrinsicHeight(width - xInset)
}

func (w *tabsElement) MinIntrinsicWidth(height base.Length) base.Length {
	xInset := w.insets.Left + w.insets.Right
	yInset := w.insets.Top + w.insets.Bottom

	if w.child == nil {
		return xInset
	}

	if height == base.Inf {
		return max(
			w.controlTabsMinWidth(),
			w.child.MinIntrinsicWidth(base.Inf)+xInset,
		)
	}

	return max(
		w.controlTabsMinWidth(),
		w.child.MinIntrinsicWidth(height-yInset),
	)
}

func (w *tabsElement) UpdateProps(data base.Widget) error {
	// Cast to correct type.
	tabs := data.(*Tabs)
	// Ensure that the Value is a useable index.
	tabs.UpdateValue()
	// Update properties.
	// Forward to the platform-dependant code where necessary.
	w.insets = tabs.Insets
	return w.updateProps(tabs)
}
