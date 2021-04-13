// +build cocoa darwin,!gtk

package goey

import (
	"guitest/goey/base"
	"guitest/goey/internal/cocoa"
)

type selectinputElement struct {
	control *cocoa.PopUpButton
}

func (w *SelectInput) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewPopUpButton(parent.Handle)
	for _, v := range w.Items {
		control.AddItem(v)
	}
	control.SetValue(w.Value, w.Unset)
	control.SetEnabled(!w.Disabled)
	control.SetCallbacks(w.OnChange, w.OnFocus, w.OnBlur)

	retval := &selectinputElement{
		control: control,
	}
	return retval, nil
}

func (w *selectinputElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *selectinputElement) Layout(bc base.Constraints) base.Size {
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *selectinputElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *selectinputElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *selectinputElement) Props() base.Widget {
	onchange, onfocus, onblur := w.control.Callbacks()

	value := w.control.Value()
	unset := value < 0
	if unset {
		value = 0
	}

	return &SelectInput{
		Items:    w.propsItems(),
		Value:    int(value),
		Unset:    unset,
		Disabled: !w.control.IsEnabled(),
		OnChange: onchange,
		OnFocus:  onfocus,
		OnBlur:   onblur,
	}
}

func (w *selectinputElement) propsItems() []string {
	length := w.control.NumberOfItems()
	ret := make([]string, length)
	for i := range ret {
		ret[i] = w.control.ItemAtIndex(i)
	}
	return ret
}

func (w *selectinputElement) TakeFocus() bool {
	return w.control.MakeFirstResponder()
}

func (w *selectinputElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *selectinputElement) updateProps(data *SelectInput) error {
	w.control.RemoveAllItems()
	for _, v := range data.Items {
		w.control.AddItem(v)
	}
	w.control.SetValue(data.Value, data.Unset)
	w.control.SetEnabled(!data.Disabled)
	w.control.SetCallbacks(data.OnChange, data.OnFocus, data.OnBlur)
	return nil
}
