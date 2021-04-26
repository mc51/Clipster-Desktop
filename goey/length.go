package goey

import (
	"clipster/goey/base"
)

// Common lengths used when describing GUIs.
// Note that the DIP (device-independent pixel) is the natural unit for this
// package.  Because of limited precision, the PT listed here is somewhat smaller
// than its correct value.
const (
	DIP  = base.DIP  // Device-independent pixel (1/96 inch)
	PT   = base.PT   // Point (1/72 inch)
	PC   = base.PC   // Pica (1/6 inch or 12 points)
	Inch = base.Inch // Inch from a British imperial system of measurements
)

func guardInf(a, b base.Length) base.Length {
	if a == base.Inf {
		return base.Inf
	}
	return b
}

func max(a, b base.Length) base.Length {
	if a > b {
		return a
	}
	return b
}

func min(a, b base.Length) base.Length {
	if a < b {
		return a
	}
	return b
}
