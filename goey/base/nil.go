package base

var (
	nilKind = NewKind("clipster/goey/base.nil")
)

// Mount will try to mount a widget.  In the case where the widget is non-nil,
// this function is a simple wrapper around calling the method Mount directly.
// If widget is nil, this function will instead return a non-nil element, but
// an element with an intrinsic size of zero and no visible elements in the
// GUI.
func Mount(parent Control, widget Widget) (Element, error) {
	if widget == nil {
		return (*nilElement)(nil), nil
	}
	return widget.Mount(parent)
}

type nilElement struct {
}

func (*nilElement) Close() {

}

func (*nilElement) Kind() *Kind {
	return &nilKind
}

func (*nilElement) Layout(bc Constraints) Size {
	if bc.IsBounded() {
		return bc.Max
	} else if bc.HasBoundedWidth() {
		return Size{bc.Max.Width, bc.Min.Height}
	} else if bc.HasBoundedHeight() {
		return Size{bc.Min.Width, bc.Max.Height}
	}
	return bc.Min
}

func (*nilElement) MinIntrinsicHeight(Length) Length {
	return 0
}

func (*nilElement) MinIntrinsicWidth(Length) Length {
	return 0
}

func (*nilElement) Props() Widget {
	return nil
}

func (*nilElement) SetBounds(Rectangle) {
	// Do nothing
}

func (*nilElement) UpdateProps(data Widget) error {
	panic("unreachable")
}
