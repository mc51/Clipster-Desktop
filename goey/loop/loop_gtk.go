// +build gtk darwin linux freebsd openbsd

package loop

import (
	"sync/atomic"
	"testing"

	"clipster/goey/internal/gtkloop"
	"clipster/goey/internal/nopanic"
)

const (
	// Flag to control behaviour of UnlockOSThread in Run.
	unlockThreadAfterRun = true
)

var (
	runLevel uint32
)

func init() {
	gtkloop.Init()
}

func initRun() error {
	// Do nothing
	return nil
}

func terminateRun() {
	// Do nothing
}

func run() {
	// Handle run level.
	if !atomic.CompareAndSwapUint32(&runLevel, 0, 1) {
		panic("internal error")
	}
	defer atomic.StoreUint32(&runLevel, 0)

	// Start the GTK loop.
	gtkloop.Run()
}

func runTesting(func() error) error {
	panic("unreachable")
}

func do(action func() error) error {
	// Make channel for the return value of the action.
	err := make(chan error, 1)

	// Depending on the run level for the main loop, either use an idle
	// callback or a higher priority callback.  The goal with using an
	// idle callback is to ensure that the system is up and running
	// before any new changes.
	if atomic.LoadUint32(&runLevel) < 2 {
		gtkloop.IdleAdd(func() {
			atomic.StoreUint32(&runLevel, 2)
			err <- nopanic.Wrap(action)
		})
	} else {
		gtkloop.MainContextInvoke(func() {
			err <- nopanic.Wrap(action)
		})
	}

	// Block on completion of action.
	return nopanic.Unwrap(<-err)
}

func stop() {
	gtkloop.Stop()
}

func testMain(m *testing.M) int {
	// On GTK, we need to be locked to a thread, but not to a particular
	// thread.  No need for special coordination.
	return m.Run()
}
