package animate

import (
	"clipster/goey/base"
)

// EaseLength encapsulates the calculations needs to smoothly vary a length.
type EaseLength struct {
	startTime  Time
	endTime    Time
	ca, cc, cd base.Length
}

func scale(v base.Length, s float64) base.Length {
	return base.Length(float64(v) * s)
}

// NewEaseLength creates a new ease to transition a length from the old value
// to the new value.  The transition time will be the same as the constant
// DefaultTransitionTime.
func NewEaseLength(newValue, oldValue base.Length) EaseLength {
	dv := newValue - oldValue
	t := CurrentTime()

	return EaseLength{
		startTime: t,
		endTime:   t + DefaultTransitionTime, /* 1 second */
		ca:        oldValue,
		cc:        3 * dv,
		cd:        -2 * dv,
	}
}

// Done returns true if the animation is completed at the specified time.
func (e *EaseLength) Done(t Time) bool {
	return t >= e.endTime
}

// Value calculate the length as a function of time.
func (e *EaseLength) Value(t Time) base.Length {
	if t >= e.endTime {
		return e.ca + e.cc + e.cd
	}

	tf := float64(t - e.startTime)
	return e.ca + scale(e.cc, tf*tf/(defaultTransitionTime*defaultTransitionTime)) +
		scale(e.cd, tf*tf*tf/(defaultTransitionTime*defaultTransitionTime*defaultTransitionTime))
}
