package goey

import (
	"clipster/goey/base"
	"testing"
)

func TestCalculateHGap(t *testing.T) {
	cases := []struct {
		w1, w2   base.Element
		expected base.Length
	}{
		{(*textinputElement)(nil), (*textinputElement)(nil), 11 * DIP}, // Space between unrelated controls
		{(*textinputElement)(nil), (*buttonElement)(nil), 11 * DIP},    // Space between unrelated controls
		{(*buttonElement)(nil), (*textinputElement)(nil), 11 * DIP},    // Space between unrelated controls
		{(*buttonElement)(nil), (*buttonElement)(nil), 7 * DIP},        // Space between adjacent buttons
	}

	for _, v := range cases {
		out := calculateHGap(v.w1, v.w2)
		if out != v.expected {
			t.Errorf("Incorrect horizontal gap calculated, %d =/= %d", out, v.expected)
		}
	}
}

func TestCalculateVGap(t *testing.T) {
	cases := []struct {
		w1, w2   base.Element
		expected base.Length
	}{
		{(*textinputElement)(nil), (*textinputElement)(nil), 11 * DIP},   // Space between unrelated controls
		{(*textinputElement)(nil), (*paragraphElement)(nil), 11 * DIP},   // Space between unrelated controls
		{(*textinputElement)(nil), (*selectinputElement)(nil), 11 * DIP}, // Space between unrelated controls
		{(*labelElement)(nil), (*textinputElement)(nil), 5 * DIP},        // Space between text labels and associated fields
		{(*labelElement)(nil), (*selectinputElement)(nil), 5 * DIP},      // Space between text labels and associated fields
		{(*labelElement)(nil), (*textareaElement)(nil), 5 * DIP},         // Space between text labels and associated fields
		{(*checkboxElement)(nil), (*checkboxElement)(nil), 7 * DIP},      // Space between related controls
		{(*paragraphElement)(nil), (*paragraphElement)(nil), 11 * DIP},   // Space between paragraphs of text
	}

	for _, v := range cases {
		out := calculateVGap(v.w1, v.w2)
		if out != v.expected {
			t.Errorf("Incorrect vertical gap calculated, %d =/= %d", out, v.expected)
		}
	}
}
