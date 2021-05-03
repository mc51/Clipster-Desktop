// +build cocoa,darwin

package loop

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"clipster/goey/internal/cocoaloop"
	"clipster/goey/internal/gtkloop"
	"clipster/goey/internal/nopanic"
	"gitlab.com/stone.code/assert"
)

const (
	// Flag to control behaviour of UnlockOSThread in Run.
	unlockThreadAfterRun = false
)

var (
	cocoaInit      sync.Once
	testingActions chan func() error
	testingSync    chan error
)

func init() {
	assert.Assert(cocoa.IsMainThread(), "Not main thread")
	runtime.LockOSThread()
}

func initRun() error {
	cocoaInit.Do(func() {
		assert.Assert(cocoa.IsMainThread(), "Not main thread")
		// cocoa.Init()
		gtkloop.Init()
	})
	return nil
}

func terminateRun() {
	// Do nothing
}

func run() {
	assert.Assert(cocoa.IsMainThread(), "Not main thread")
	// cocoa.Run()
	// Start the GTK loop.
	gtkloop.Run()
}

func runTesting(action func() error) error {
	testingActions <- action
	return nopanic.Unwrap(<-testingSync)
}

func do(action func() error) error {
	return cocoa.PerformOnMainThread(action)
}

func stop() {
	// cocoa.Stop()
	gtkloop.Stop()
}

func testMain(m *testing.M) int {
	// Ensure that we are locked to the main thread.
	runtime.LockOSThread()
	assert.Assert(cocoa.IsMainThread(), "Not main thread")

	atomic.StoreUint32(&isTesting, 1)
	defer func() {
		atomic.StoreUint32(&isTesting, 0)
	}()

	testingActions = make(chan func() error)
	testingSync = make(chan error)

	// call flag.Parse() here if TestMain uses flags
	wait := make(chan int, 1)
	go func() {
		wait <- m.Run()
		close(testingActions)
	}()

	for a := range testingActions {
		assert.Assert(cocoa.IsMainThread(), "Not main thread")

		err := func() (err error) {
			atomic.StoreUint32(&isTesting, 0)
			defer func() {
				atomic.StoreUint32(&isTesting, 1)
			}()

			return nopanic.Wrap(func() error {
				return Run(a)
			})
		}()
		testingSync <- err
	}

	return <-wait
}
