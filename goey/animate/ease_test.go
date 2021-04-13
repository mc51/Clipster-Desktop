package animate

import (
	"guitest/goey/base"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func easeValues(values []reflect.Value, rand *rand.Rand) {
	// Create a choices for the limits of the transition.
	// Cannot use the full range of Int64 for the length, as this will lead
	// to overflow during the calculations.
	values[0] = reflect.ValueOf(base.Length(rand.Int63()) / 6)
	values[1] = reflect.ValueOf(base.Length(rand.Int63()) / 6)
}

func TestEaseLength(t *testing.T) {
	cases := []struct {
		startValue, endValue base.Length
		time                 Time
		out                  base.Length
		done                 bool
	}{
		{100 * base.DIP, 0, 0, 0, false},
		{100 * base.DIP, 0, 500, 50 * base.DIP, false},
		{100 * base.DIP, 0, 1000, 100 * base.DIP, true},
		{100 * base.DIP, 0, 1500, 100 * base.DIP, true},
	}

	for i, v := range cases {
		ease := NewEaseLength(v.startValue, v.endValue)
		if out := ease.Value(ease.startTime + v.time); out != v.out {
			t.Errorf("Case %d:  Incorrect value, want %v, got %v", i, v.out, out)
		}
		if done := ease.Done(ease.startTime + v.time); done != v.done {
			t.Errorf("Case %d:  Incorrect value, want %v, got %v", i, v.done, done)
		}
	}

	t.Run("Quick", func(t *testing.T) {
		f := func(newValue, oldValue base.Length) bool {
			ease := NewEaseLength(newValue, oldValue)
			t1 := ease.startTime
			t2 := ease.endTime
			tmp := ease.Value((t1+t2)/2) - (newValue+oldValue)/2
			return ease.Value(t1) == oldValue && ease.Value(t2) == newValue &&
				tmp >= -1 && tmp <= 1 &&
				!ease.Done(t2-1) && ease.Done(t2)
		}
		if err := quick.Check(f, &quick.Config{Values: easeValues}); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}
