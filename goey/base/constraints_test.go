package base

import (
	"testing"
)

func TestConstraints(t *testing.T) {
	cases := []struct {
		in                                                   Constraints
		isNormalized, isTight, hasTightWidth, hasTightHeight bool
		isBounded, hasBoundedWidth, hasBoundedHeight         bool
	}{
		{Expand(), true, false, false, false, false, false, false},
		{ExpandHeight(10 * DIP), true, false, true, false, false, true, false},
		{ExpandWidth(10 * DIP), true, false, false, true, false, false, true},
		{Loose(Size{10 * DIP, 15 * DIP}), true, false, false, false, true, true, true},
		{Tight(Size{10 * DIP, 15 * DIP}), true, true, true, true, true, true, true},
		{TightWidth(10 * DIP), true, false, true, false, false, true, false},
		{TightHeight(10 * DIP), true, false, false, true, false, false, true},
	}

	for i, v := range cases {
		if out := v.in.IsNormalized(); v.isNormalized != out {
			t.Errorf("Failed on case %d for IsNormalized, want %v, got %v", i, v.isNormalized, out)
		}
		if out := v.in.IsTight(); v.isTight != out {
			t.Errorf("Failed on case %d for IsTight, want %v, got %v", i, v.isTight, out)
		}
		if out := v.in.HasTightWidth(); v.hasTightWidth != out {
			t.Errorf("Failed on case %d for HasTightWidth, want %v, got %v", i, v.hasTightWidth, out)
		}
		if out := v.in.HasTightHeight(); v.hasTightHeight != out {
			t.Errorf("Failed on case %d for HasTightHeight, want %v, got %v", i, v.hasTightHeight, out)
		}
		if out := v.in.IsBounded(); v.isBounded != out {
			t.Errorf("Failed on case %d for IsBounded, want %v, got %v", i, v.isBounded, out)
		}
		if out := v.in.HasBoundedWidth(); v.hasBoundedWidth != out {
			t.Errorf("Failed on case %d for HasBoundedWidth, want %v, got %v", i, v.hasBoundedWidth, out)
		}
		if out := v.in.HasBoundedHeight(); v.hasBoundedHeight != out {
			t.Errorf("Failed on case %d for HasBoundedHeight, want %v, got %v", i, v.hasBoundedHeight, out)
		}

		if out := v.in.IsZero(); out {
			t.Errorf("Failed on case %d for IsZero, want %v, got %v", i, false, out)
		}
	}
}

func TestConstraints_Constrain(t *testing.T) {
	cases := []struct {
		in   Constraints
		size Size
		out  Size
	}{
		{Tight(Size{10 * DIP, 15 * DIP}), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Tight(Size{10 * DIP, 15 * DIP}), Size{20 * DIP, 30 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Tight(Size{10 * DIP, 15 * DIP}), Size{2 * DIP, 3 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightWidth(10 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightWidth(15 * DIP), Size{10 * DIP, 15 * DIP}, Size{15 * DIP, 15 * DIP}},
		{TightWidth(5 * DIP), Size{10 * DIP, 15 * DIP}, Size{5 * DIP, 15 * DIP}},
		{TightHeight(15 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightHeight(30 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 30 * DIP}},
		{TightHeight(75 * DIP / 10), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 75 * DIP / 10}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{20 * DIP, 30 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{2 * DIP, 3 * DIP}, Size{2 * DIP, 3 * DIP}},
	}

	for i, v := range cases {
		if out := v.in.Constrain(v.size); v.out != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
		if out := v.in.ConstrainWidth(v.size.Width); v.out.Width != out {
			t.Errorf("Failed on case %d width, want %v, got %v", i, v.out.Width, out)
		}
		if out := v.in.ConstrainHeight(v.size.Height); v.out.Height != out {
			t.Errorf("Failed on case %d height, want %v, got %v", i, v.out.Height, out)
		}
	}
}

func TestConstraints_ConstrainAndAttemptToPreserveAspectRatio(t *testing.T) {
	cases := []struct {
		in   Constraints
		size Size
		out  Size
	}{
		{Tight(Size{10 * DIP, 15 * DIP}), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Tight(Size{10 * DIP, 15 * DIP}), Size{20 * DIP, 30 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Tight(Size{10 * DIP, 15 * DIP}), Size{2 * DIP, 3 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightWidth(10 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightWidth(15 * DIP), Size{10 * DIP, 15 * DIP}, Size{15 * DIP, 225 * DIP / 10}},
		{TightWidth(5 * DIP), Size{10 * DIP, 15 * DIP}, Size{5 * DIP, 75 * DIP / 10}},
		{TightHeight(15 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightHeight(30 * DIP), Size{10 * DIP, 15 * DIP}, Size{20 * DIP, 30 * DIP}},
		{TightHeight(75 * DIP / 10), Size{10 * DIP, 15 * DIP}, Size{5 * DIP, 75 * DIP / 10}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{20 * DIP, 30 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{2 * DIP, 3 * DIP}, Size{2 * DIP, 3 * DIP}},
	}

	for i, v := range cases {
		out := v.in.ConstrainAndAttemptToPreserveAspectRatio(v.size)
		if v.out != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
	}
}
func TestConstraints_Inset(t *testing.T) {
	cases := []struct {
		in      Constraints
		deflate Length
		out     Constraints
	}{
		{Tight(Size{}), 1 * DIP, Tight(Size{})},
		{Tight(Size{2 * DIP, 2 * DIP}), 10 * DIP, Tight(Size{})},
		{Tight(Size{10 * DIP, 11 * DIP}), 5 * DIP, Tight(Size{5 * DIP, 6 * DIP})},
		{Loose(Size{}), 1 * DIP, Loose(Size{})},
		{Loose(Size{2 * DIP, 2 * DIP}), 10 * DIP, Loose(Size{})},
		{Loose(Size{10 * DIP, 11 * DIP}), 5 * DIP, Loose(Size{5 * DIP, 6 * DIP})},
		{TightWidth(0), 1 * DIP, TightWidth(0)},
		{TightWidth(2 * DIP), 10 * DIP, TightWidth(0)},
		{TightWidth(10 * DIP), 5 * DIP, TightWidth(5 * DIP)},
		{TightHeight(0), 1 * DIP, TightHeight(0)},
		{TightHeight(2 * DIP), 10 * DIP, TightHeight(0)},
		{TightHeight(10 * DIP), 5 * DIP, TightHeight(5 * DIP)},
		{Expand(), 5 * DIP, Expand()},
	}

	for i, v := range cases {
		out := v.in.Inset(v.deflate, v.deflate)
		if v.out != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
	}
}

func TestConstraints_Loosen(t *testing.T) {
	size := Size{10 * DIP, 15 * DIP}
	sizeZ := Size{}

	cases := []struct {
		in        Constraints
		out       Constraints
		outWidth  Constraints
		outHeight Constraints
	}{
		{Tight(size), Constraints{sizeZ, size}, Constraints{Size{0, 15 * DIP}, size}, Constraints{Size{10 * DIP, 0}, size}},
		{Loose(size), Constraints{sizeZ, size}, Constraints{sizeZ, size}, Constraints{sizeZ, size}},
	}

	for i, v := range cases {
		if out := v.in.Loosen(); v.out != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		} else if out.HasTightHeight() {
			t.Errorf("Failed on case %d, has tight height", i)
		} else if out.HasTightWidth() {
			t.Errorf("Failed on case %d, has tight width", i)
		}

		if out := v.in.LoosenHeight(); v.outHeight != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.outHeight, out)
		} else if out.HasTightHeight() {
			t.Errorf("Failed on case %d, has tight height", i)
		}

		if out := v.in.LoosenWidth(); v.outWidth != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.outWidth, out)
		} else if out.HasTightWidth() {
			t.Errorf("Failed on case %d, has tight height", i)
		}
	}
}

func TestConstraints_IsSatisfiedBy(t *testing.T) {
	defSize := Size{10 * DIP, 15 * DIP}

	cases := []struct {
		in   Constraints
		size Size
		out  bool
	}{
		{Loose(defSize), defSize, true},
		{Loose(defSize), Size{}, true},
		{Loose(defSize), Size{Height: defSize.Height}, true},
		{Loose(defSize), Size{Width: defSize.Width}, true},
		{Loose(defSize), Size{defSize.Width + 1, defSize.Height + 1}, false},
		{Tight(defSize), defSize, true},
		{Tight(defSize), Size{}, false},
		{Tight(defSize), Size{Height: defSize.Height}, false},
		{Tight(defSize), Size{Width: defSize.Width}, false},
		{Tight(defSize), Size{defSize.Width + 1, defSize.Height + 1}, false},
	}

	for i, v := range cases {
		if out := v.in.IsSatisfiedBy(v.size); v.out != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
	}
}

func TestConstraints_Tighten(t *testing.T) {
	size1 := Size{10 * DIP, 10 * DIP}

	cases := []struct {
		in   Constraints
		size Size
		out  Constraints
		outH Constraints
		outV Constraints
	}{
		{Expand(), size1, Tight(size1), ExpandWidth(size1.Height), ExpandHeight(size1.Width)},
		{ExpandHeight(10 * DIP), size1, Tight(size1), Tight(size1), ExpandHeight(10 * DIP)},
		{ExpandWidth(10 * DIP), size1, Tight(size1), ExpandWidth(10 * DIP), Tight(size1)},
		{Loose(Size{20 * DIP, 25 * DIP}), size1, Tight(size1), Constraints{Size{0, 10 * DIP}, Size{20 * DIP, 10 * DIP}}, Constraints{Size{10 * DIP, 0}, Size{10 * DIP, 25 * DIP}}},
		{Loose(Size{20 * DIP, 25 * DIP}), Size{30 * DIP, 30 * DIP}, Tight(Size{20 * DIP, 25 * DIP}), Constraints{Size{0, 25 * DIP}, Size{20 * DIP, 25 * DIP}}, Constraints{Size{20 * DIP, 0}, Size{20 * DIP, 25 * DIP}}},
		{Tight(Size{10 * DIP, 15 * DIP}), size1, Tight(Size{10 * DIP, 15 * DIP}), Tight(Size{10 * DIP, 15 * DIP}), Tight(Size{10 * DIP, 15 * DIP})},
		{TightWidth(15 * DIP), size1, Tight(Size{15 * DIP, 10 * DIP}), Tight(Size{15 * DIP, 10 * DIP}), TightWidth(15 * DIP)},
		{TightHeight(15 * DIP), size1, Tight(Size{10 * DIP, 15 * DIP}), TightHeight(15 * DIP), Tight(Size{10 * DIP, 15 * DIP})},
	}

	for i, v := range cases {
		if out := v.in.Tighten(v.size); out != v.out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
		if out := v.in.TightenHeight(v.size.Height); out != v.outH {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.outH, out)
		}
		if out := v.in.TightenWidth(v.size.Width); out != v.outV {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.outV, out)
		}
	}
}
