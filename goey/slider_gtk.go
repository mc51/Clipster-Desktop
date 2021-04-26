// +build gtk linux,!cocoa freebsd,!cocoa openbsd,!cocoa

package goey

import (
	"clipster/goey/base"
	"clipster/goey/internal/gtk"
)

type sliderElement struct {
	Control
	value    float64
	min, max float64

	onChange func(float64)
	onFocus  func()
	onBlur   func()
}

func (w *Slider) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountSlider(parent.Handle, w.Value, w.Disabled, w.Min, w.Max,
		w.OnChange != nil, w.OnFocus != nil, w.OnBlur != nil)

	retval := &sliderElement{
		Control:  Control{control},
		value:    w.Value,
		min:      w.Min,
		max:      w.Max,
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}

func (w *sliderElement) OnChange(value float64) {
	if value != w.value {
		w.value = value
		w.onChange(value)
	}
}

func (w *sliderElement) OnFocus() {
	w.onFocus()
}

func (w *sliderElement) OnBlur() {
	w.onBlur()
}

func (w *sliderElement) Props() base.Widget {
	return &Slider{
		Value:    gtk.SliderValue(w.handle),
		Min:      w.min,
		Max:      w.max,
		Disabled: !gtk.WidgetSensitive(w.handle),
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

// Layout determines the best size for an element that satisfies the
// constraints.
func (w *sliderElement) Layout(bc base.Constraints) base.Size {
	if !bc.HasBoundedWidth() && !bc.HasBoundedHeight() {
		// No need to worry about breaking the constraints.  We can take as
		// much space as desired.
		width := w.MinIntrinsicWidth(base.Inf)
		height := gtk.WidgetNaturalHeight(w.handle)
		// Dimensions may need to be increased to meet minimums.
		return bc.Constrain(base.Size{width, base.FromPixelsY(height)})
	}
	if !bc.HasBoundedHeight() {
		// No need to worry about height.  Find the width that best meets the
		// widgets preferred width.
		width := bc.ConstrainWidth(w.MinIntrinsicWidth(base.Inf))
		// Get the best height for this width.
		height := gtk.WidgetNaturalHeightForWidth(w.handle, width.PixelsX())
		// Height may need to be increased to meet minimum.
		return base.Size{width, bc.ConstrainHeight(base.FromPixelsY(height))}
	}

	// Not clear the following is the best general approach given GTK layout
	// model.
	height2 := gtk.WidgetNaturalHeight(w.handle)
	if height := base.FromPixelsY(height2); height < bc.Max.Height {
		width := w.MinIntrinsicWidth(height)
		return bc.Constrain(base.Size{width, height})
	}

	height := base.FromPixelsY(height2)
	width := w.MinIntrinsicWidth(height)
	return bc.Constrain(base.Size{width, height})
}

// MinIntrinsicWidth returns the minimum width that this element requires
// to be correctly displayed.
func (w *sliderElement) MinIntrinsicWidth(base.Length) base.Length {
	width := gtk.WidgetMinWidth(w.handle)
	if limit := base.FromPixelsX(width); limit < 160*DIP {
		return 160 * DIP
	}
	return base.FromPixelsX(width)
}

func (w *sliderElement) updateProps(data *Slider) error {
	gtk.SliderUpdate(w.handle, data.Value, data.Disabled, data.Min, data.Max,
		data.OnChange != nil, data.OnFocus != nil, data.OnBlur != nil)

	w.min = data.Min
	w.max = data.Max
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	return nil
}
