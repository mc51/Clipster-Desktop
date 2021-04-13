package animate

import (
	"guitest/goey/loop"
	"errors"
	"sync/atomic"
	"time"
)

// Time represents an instant in time with millisecond precision.
type Time uint

const (
	defaultTransitionTime = 1000
	// The default transition time is the default time for an animation.
	DefaultTransitionTime Time = defaultTransitionTime
)

// An element that wishes to be animated must implement this interface.
type Element interface {
	// AnimateFrame will be called for each frame to provide the element the
	// opportunity to modify its state.  The method should return true to
	// continue to receive callbacks.  Once the element's animations are
	// completed, it should return false.
	AnimateFrame(time Time) bool
}

var (
	currentTime      Time
	isRunning        uint32
	elements         = map[Element]struct{}{}
	errStopAnimation = errors.New("stop animation")
)

// CurrentTime returns the time for the current animation frame.  If no
// animations are currently running, it will return the current local
// time.
func CurrentTime() Time {
	if atomic.LoadUint32(&isRunning) != 0 {
		return currentTime
	}

	return Time(time.Now().UnixNano() / 1e6)
}

// AddAnimation adds the element to the set of elements with active animations.
// As necessary, this function will arrange for the element to be updated with
// the time of the current animation frame so that its state can be updated.
//
// This function should only be called from the GUI thread.
func AddAnimation(elem Element) {
	// Access the map of running animations should only happen on the GUI thread.
	// No synchronization is required.
	// Add the element to the set of running animations.
	elements[elem] = struct{}{}

	if atomic.CompareAndSwapUint32(&isRunning, 0, 1) {
		go animateFrames()
	}
}

func animateFrames() {
	defer atomic.StoreUint32(&isRunning, 0)

	// TODO:  Adjust frame rate if we cannot maintain 60 frames per second
	ticker := time.NewTicker(time.Second / 30)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case t := <-ticker.C:
			// Update the time for this frame.
			currentTime = Time(t.UnixNano() / 1e6)
			// Callback to the GUI thread to update the elements.
			err := loop.Do(animateFrame)
			if err == errStopAnimation {
				return
			}
		}
	}
}

// To animate a frame, we need to notify all participating elements of the
// current time.
func animateFrame() error {
	for k := range elements {
		// Update state of the element.
		ok := k.AnimateFrame(currentTime)
		if !ok {
			// No more animation for this element.
			// Remove from set.
			delete(elements, k)
		}
	}

	// If there is nothing to animate, we can close the goroutine.
	if len(elements) == 0 {
		return errStopAnimation
	}

	return nil
}
