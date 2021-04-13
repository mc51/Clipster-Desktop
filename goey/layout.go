package goey

import (
	"guitest/goey/base"
)

func calculateHGap(previous base.Element, current base.Element) base.Length {
	// The vertical gap between most controls is 11 relative pixels.  However,
	// there are different rules for between a label and its associated control,
	// or between related controls.  These relationship do not appear in the
	// model provided by this package, so these relationships need to be
	// inferred from the order and type of controls.
	//
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	if _, ok := previous.(*buttonElement); ok {
		if _, ok := current.(*buttonElement); ok {
			// Any pair of successive buttons will be assumed to be in a
			// related group.
			return 7 * DIP
		}
	}

	// The spacing between unrelated controls.
	return 11 * DIP
}

func calculateVGap(previous base.Element, current base.Element) base.Length {
	// The vertical gap between most controls is 11 relative pixels.  However,
	// there are different rules for between a label and its associated control,
	// or between related controls.  These relationship do not appear in the
	// model provided by this package, so these relationships need to be
	// inferred from the order and type of controls.
	//
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing

	// Unwrap and Expand widgets.
	if expand, ok := previous.(*expandElement); ok {
		previous = expand.child
	}
	if expand, ok := current.(*expandElement); ok {
		current = expand.child
	}

	// Apply layout rules.
	if _, ok := previous.(*labelElement); ok {
		// Any label immediately preceding any other control will be assumed to
		// be 'associated'.
		return 5 * DIP
	}
	if _, ok := previous.(*checkboxElement); ok {
		if _, ok := current.(*checkboxElement); ok {
			// Any pair of successive checkboxes will be assumed to be in a
			// related group.
			return 7 * DIP
		}
	}

	// The spacing between unrelated controls.  This is also the default space
	// between paragraphs of text.
	return 11 * DIP
}
