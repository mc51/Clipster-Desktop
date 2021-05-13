// +build gtk linux darwin freebsd openbsd

package goey

import (
	"clipster/goey/base"
	"clipster/goey/internal/gtk"
)

type textareaElement struct {
	Control

	minLines int
	onChange func(string)
	onFocus  func()
	onBlur   func()
}

func (w *TextArea) mount(parent base.Control) (base.Element, error) {
	control := gtk.MountTextarea(parent.Handle, w.Value, w.Disabled, w.ReadOnly,
		w.OnChange != nil, w.OnFocus != nil, w.OnBlur != nil)

	retval := &textareaElement{
		Control:  Control{control},
		minLines: minlinesDefault(w.MinLines),
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	}
	gtk.RegisterWidget(control, retval)

	return retval, nil
}

func (w *textareaElement) Layout(bc base.Constraints) base.Size {
	if !bc.HasBoundedWidth() {
		if bc.Min.Width > 0 {
			width := bc.Min.Width
			height := w.MinIntrinsicHeight(width)
			return bc.Constrain(base.Size{width, height})
		}

		width := gtk.WidgetNaturalWidth(w.handle)
		height := w.MinIntrinsicHeight(base.Inf)
		return bc.Constrain(base.Size{
			base.FromPixelsX(width),
			height,
		})
	}

	width := bc.Max.Width
	height := w.MinIntrinsicHeight(width)
	return bc.Constrain(base.Size{width, height})
}

func (w *textareaElement) MinIntrinsicHeight(width base.Length) base.Length {
	// This won't respond correctly to changes in font size on GTK, but
	// we need to establish a height to set minlines.
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	const lineHeight = 16 * DIP
	minHeight := 23*DIP + lineHeight.Scale(w.minLines-1, 1)

	if width != base.Inf {
		height := gtk.WidgetMinHeightForWidth(w.handle, width.PixelsX())
		return max(minHeight, base.FromPixelsY(height))
	}
	height := gtk.WidgetMinHeight(w.handle)
	return max(minHeight, base.FromPixelsY(height))
}

func (w *textareaElement) MinIntrinsicWidth(base.Length) base.Length {
	width := gtk.WidgetMinWidth(w.handle)
	return base.FromPixelsX(width)
}

func (w *textareaElement) OnChange(value string) {
	if w.onChange != nil {
		w.onChange(value)
	}
}

func (w *textareaElement) OnFocus() {
	w.onFocus()
}

func (w *textareaElement) OnBlur() {
	w.onBlur()
}

func (w *textareaElement) OnEnterKey(value string) {
	// Not supported by GTK.
	// This event will never occur
}

func (w *textareaElement) Props() base.Widget {
	return &TextArea{
		Value:    gtk.TextareaText(w.handle),
		Disabled: !gtk.WidgetSensitive(gtk.TextareaTextview(w.handle)),
		ReadOnly: gtk.TextareaReadOnly(w.handle),
		MinLines: w.minLines,
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *textareaElement) TakeFocus() bool {
	control := Control{gtk.TextareaTextview(w.handle)}
	return control.TakeFocus()
}

func (w *textareaElement) TypeKeys(text string) chan error {
	control := Control{gtk.TextareaTextview(w.handle)}
	return control.TypeKeys(text)
}

func (w *textareaElement) updateProps(data *TextArea) error {
	w.onChange = nil // temporarily break OnChange to prevent event
	gtk.TextareaUpdate(w.handle, data.Value, data.Disabled, data.ReadOnly,
		data.OnChange != nil, data.OnFocus != nil, data.OnBlur != nil)

	w.minLines = data.MinLines
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}
